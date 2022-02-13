package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"development/mxshop_api/user-web/forms"
	"development/mxshop_api/user-web/global"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func GenerateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().Unix())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", rand.Intn(r))
		if err != nil {
			return ""
		}
	}

	return sb.String()
}

func SendSms(c *gin.Context) {
	smsCode := GenerateSmsCode(6)

	sendSmsForm := forms.SendSmsForm{}
	err := c.ShouldBind(&sendSmsForm)
	if err != nil {
		HandleValidatorError(c, err)
		return
	}

	client, err := dysmsapi.NewClientWithAccessKey("cn-beijing", global.SrvConfig.AliSmsInfo.ApiKey, global.SrvConfig.AliSmsInfo.ApiSecret)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = sendSmsForm.Mobile            //手机号
	request.QueryParams["SignName"] = "慕学在线"                            //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = "SMS_234158408"               //阿里云的短信模板号 自己设置
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //短信模板中的验证码内容 自己生成   之前试过直接返回，但是失败，加上code成功。

	response, err := client.ProcessCommonRequest(request)

	fmt.Println(client.DoAction(request, response))
	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"sms": "发送验证码失败",
		})
		return
	}

	//将验证码保存起来 - redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.SrvConfig.RedisInfo.Host, global.SrvConfig.RedisInfo.Port),
	})
	rdb.Set(sendSmsForm.Mobile, smsCode, time.Duration(global.SrvConfig.RedisInfo.Expire)*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
