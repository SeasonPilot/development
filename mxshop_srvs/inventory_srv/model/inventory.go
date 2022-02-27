package model

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index"` // 经常查询的要加索引
	Stocks  int32 `gorm:"type:int"`
	Version int32 `gorm:"type:int"` //分布式锁的乐观锁
}
