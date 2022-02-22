package response

import (
	"fmt"
	"time"
)

type GoodsRsp struct {
	Id       uint32   `json:"id"`
	BirthDay JsonTime `json:"birthday"`
}

type JsonTime time.Time

// MarshalJSON 接受者不是 (u *GoodsRsp)，是字段的类型
// 转换 BirthDay 格式
func (t JsonTime) MarshalJSON() ([]byte, error) {
	// 类型转换，将 t 转换成 time.Time 类型
	birthday := time.Time(t).Format("2006-01-02")

	// 不加这个步骤就会报错  nvalid character '-' after top-level value 为什么？？？？
	// 就加了个双引号 "2022-02-08"
	s := fmt.Sprintf("\"%s\"", birthday)

	return []byte(s), nil
}
