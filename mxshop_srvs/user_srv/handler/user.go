package handler

import (
	"context"

	"development/mxshop_srvs/user_srv/global"
	"development/mxshop_srvs/user_srv/model"
	"development/mxshop_srvs/user_srv/proto"

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
