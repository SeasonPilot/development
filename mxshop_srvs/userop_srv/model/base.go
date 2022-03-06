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
	ID        int32          `gorm:"type:int" json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
	IsDelete  bool           `json:"-"`
}
