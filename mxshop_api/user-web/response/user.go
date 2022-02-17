package response

import (
	"fmt"
	"time"
)

type UserRsp struct {
	Id uint32 `json:"id"`
	//Password string   `json:"password"`
	Mobile   string   `json:"mobile"`
	NickName string   `json:"name"`
	BirthDay JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Role     uint32   `json:"role"`
}

type JsonTime time.Time

// MarshalJSON 接受者不是 (u *UserRsp)，是字段的类型
// 转换 BirthDay 格式
func (t JsonTime) MarshalJSON() ([]byte, error) {
	// 类型转换，将 t 转换成 time.Time 类型
	birthday := time.Time(t).Format("2006-01-02")

	// 不加这个步骤就会报错  nvalid character '-' after top-level value 为什么？？？？
	// 就加了个双引号 "2022-02-08"
	s := fmt.Sprintf("\"%s\"", birthday)

	return []byte(s), nil
}
