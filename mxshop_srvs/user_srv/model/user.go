package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID       int        `gorm:"primarykey"`
	CreateAt *time.Time `gorm:"column:add_time"`
	UpdateAt *time.Time `gorm:"column:update_time"`
	DeleteAt gorm.DeletedAt
	isDelete bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"type:varchar(11);not null;unique;index:idx_mobile"`
	Password string     `gorm:"type:varchar(100);not null"`
	NickName string     `gorm:"type:varchar(20)"`
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"type:varchar(6);default:male comment 'male表示男性,female表示女性'"`
	Role     int        `gorm:"default:1 comment '1 表示普通用户,2 表示管理员'"`
}
