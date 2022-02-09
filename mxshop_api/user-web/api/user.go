package api

import (
	"fmt"
	"net/http"
	"time"

	"development/mxshop_api/user-web/response"
	"development/mxshop_srvs/user_srv/proto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

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

func GetUserList(c *gin.Context) {
	ip := "127.0.0.1"
	port := 50051
	// non-blocking dial
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicf("连接用户服务失败 %s", err.Error())
		return
	}

	userClient := proto.NewUserClient(conn)
	rsp, err := userClient.GetUserList(c, &proto.PageInfo{
		Pn:    0,
		PSize: 10,
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
