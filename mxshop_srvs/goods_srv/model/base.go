package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type GormList []string

func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

type BaseModel struct {
	ID        int32     `gorm:"type:int"`
	CreatedAt time.Time `gorm:"column:create_time"`
	UpdateAt  time.Time `gorm:"column:update_time"`
	DeleteAt  gorm.DeletedAt
	IsDelete  bool
}
