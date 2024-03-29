package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
removeTopStruct 移除 `PassWordLoginForm.`
{
    "msg": {
        "PassWordLoginForm.password": "password长度必须至少为3个字符"
    }
}
*/
func removeTopStruct(fields map[string]string) map[string]string {
	resp := make(map[string]string, len(fields))
	for k, v := range fields {
		resp[k[strings.Index(k, ".")+1:]] = v
	}
	return resp
}

// RpcErrToHttpErr 将 grpc 的 code 转换成 http 的状态码
func RpcErrToHttpErr(err error, c *gin.Context) {
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			switch grpcStatus.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": grpcStatus.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					// 不能直接返回  grpcStatus.Message() ,会暴露过多信息给用户,如敏感信息。
					// 不能把 grpc 内部错误暴露给用户，不友好
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": grpcStatus.Code(),
				})
			}
		}
	}
}

// HandleValidatorError 处理 Validator 的错误
func HandleValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": removeTopStruct(errs.Translate(global.Translator)),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": err.Error(),
	})
}

func GetUserList(c *gin.Context) {
	// 从 Context 中获取当前登陆的用户
	claims, _ := c.Get("claims")
	currentUser, _ := claims.(*models.CustomClaims)
	zap.S().Infof("当前登陆的用户是: %d", currentUser.ID)

	// 解析请求参数
	pn := c.DefaultQuery("pn", "1")
	pnInt, _ := strconv.Atoi(pn)
	pSize := c.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserClient.GetUserList(c, &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorf("GetUserList 失败 %s", err.Error())
		RpcErrToHttpErr(err, c)
		return
	}

	// make 的第二个参数不能是 len(rsp.Data)，这样会先创建一个长度为 len(rsp.Data) 的零值的 slice
	// - slice 初始化预先分配内存可以提升性能；直接使用 index 而非 append 可以提升性能；
	users := make([]response.UserRsp, len(rsp.Data))
	for i, data := range rsp.Data {
		user := response.UserRsp{
			Id: data.Id,
			//Password: data.Password,
			Mobile:   data.Mobile,
			NickName: data.NickName,
			BirthDay: response.JsonTime(time.Unix(int64(data.BirthDay), 0)),
			Gender:   data.Gender,
			Role:     data.Role,
		}
		users[i] = user
	}

	c.JSON(http.StatusOK, users)
}

func PassWordLogin(c *gin.Context) {
	// 1.表单验证
	passwordLoginForm := forms.PassWordLoginForm{}
	err := c.ShouldBind(&passwordLoginForm)
	if err != nil {
		HandleValidatorError(c, err)
		// fixme: 要记的 return
		return
	}

	// 2.Captcha 验证
	if !store.Verify(passwordLoginForm.CaptchaID, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "验证码错误",
		})
		return
	}

	// 3.查看用户是否存在
	rsp, err := global.UserClient.GetUserByMobile(c, &proto.MobileRequest{Mobile: passwordLoginForm.Name})
	if err != nil {
		// fixme: user_srv grpc 返回的错误不只一种,可以看下 grpc 层服务返回哪些错误; 所以这里要拿到错误原因进行判断。
		//  还有连接不上 grpc 服务的错误 Unavailable。
		RpcStatus, ok := status.FromError(err)
		zap.S().Errorf("rpcErrMsg: %s", RpcStatus.Message())
		if ok {
			switch RpcStatus.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户不存在",
				})
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "rpc 服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登陆失败",
				})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登陆失败",
		})
		zap.S().Errorf("GetUserByMobile err:%s", err.Error())
		return
	}
	if rsp == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "用户不存在",
		})
		return
	}

	// 4.验证密码
	checkPassWordRsp, err := global.UserClient.CheckPassWord(c, &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: rsp.Password,
	})
	if err != nil {
		RpcStatus, ok := status.FromError(err)
		if ok {
			switch RpcStatus.Code() {
			case codes.Unavailable:
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"msg": "rpc 服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "登陆失败",
				})
			}
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登陆失败",
		})
		return
	}

	if !checkPassWordRsp.Success {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "密码错误",
		})
		return
	}

	// 5.登陆成功后返回 JWT Token
	j := middlewares.NewJWT()
	token, err := j.CreateToken(
		// 注意，不要在 JWT 的 payload 或 header 中放置敏感信息，除非它们是加密的
		models.CustomClaims{
			ID:          uint(rsp.Id),
			NickName:    rsp.NickName,
			AuthorityID: uint(rsp.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "season",
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)}, //30天过期
				NotBefore: &jwt.NumericDate{Time: time.Now()},
			},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	// 是客户端请求时需要带上 x-token ,服务端通过 body 返回，不是 response header
	//c.Header("x-token", token)
	c.JSON(http.StatusOK, gin.H{
		"id":         rsp.Id,
		"nick_name":  rsp.NickName,
		"token":      token,
		"expired_at": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
}

func Register(c *gin.Context) {
	// 1.表单验证
	var registerForm forms.RegisterForm
	err := c.ShouldBind(&registerForm)
	if err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 2. sms 验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.SrvConfig.RedisInfo.Host, global.SrvConfig.RedisInfo.Port),
	})
	result, err := rdb.Get(registerForm.Mobile).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "短信验证码错误",
		})
		return
	}

	if registerForm.Code != result {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "短信验证码错误",
		})
		return
	}

	// 3.注册用户
	user, err := global.UserClient.CreateUser(c, &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		RpcErrToHttpErr(err, c)
		return
	}

	// 4.注册成功后返回 JWT Token, 实现 注册即登录; 注册后自动登录。
	j := middlewares.NewJWT()
	token, err := j.CreateToken(
		// 注意，不要在 JWT 的 payload 或 header 中放置敏感信息，除非它们是加密的
		models.CustomClaims{
			ID:          uint(user.Id),
			NickName:    user.NickName,
			AuthorityID: uint(user.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "season",
				ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 30)}, //30天过期
				NotBefore: &jwt.NumericDate{Time: time.Now()},
			},
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
}
