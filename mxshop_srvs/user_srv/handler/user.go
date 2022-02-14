package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"mxshop-srvs/user_srv/global"
	"mxshop-srvs/user_srv/model"
	"mxshop-srvs/user_srv/proto"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserServer struct{}

func modelToResp(user model.User) *proto.UserInfoResponse {
	//在 grpc 的 message 中字段有默认值，你不能随便赋值 nil 进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoResp := &proto.UserInfoResponse{
		Id:       uint32(user.ID),
		Password: user.Password,
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     uint32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.BirthDay = uint32(user.Birthday.Unix())
	}
	return userInfoResp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context, in *proto.PageInfo) (resp *proto.UserListResponse, err error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	// fixme：指针类型要先初始化再使用
	resp = &proto.UserListResponse{}
	resp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(in.Pn), int(in.PSize))).Find(&users)
	for _, user := range users {
		resp.Data = append(resp.Data, modelToResp(user))
	}

	return resp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	user := model.User{}
	result := global.DB.Where(model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return modelToResp(user), nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	user := model.User{}
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return modelToResp(user), nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 先查询用户是否已经创建
	var user model.User
	result := global.DB.Where(model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Error(codes.AlreadyExists, "用户已存在")
	}

	user.NickName = req.NickName
	user.Mobile = req.Mobile
	//密码加密
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return modelToResp(user), nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	// 查询用户是否存在
	var user model.User
	result := global.DB.First(&user, req.ID)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	// 更新字段
	user.NickName = req.NickName
	user.Gender = req.Gender
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.Birthday = &birthDay

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新失败:%s", result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	return &proto.CheckResponse{
		Success: password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options),
	}, nil
}
