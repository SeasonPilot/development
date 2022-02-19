package model

type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(50);not null;" json:"name,omitempty"`
	ParentCategoryID int32       `gorm:"type:int" json:"parent_category_id,omitempty"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;reference:ID" json:"sub_category,omitempty"`
	Level            int32       `gorm:"type:int;default:1;not null" json:"level,omitempty"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab,omitempty"`
}

type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(50);not null"`
	Logo string `gorm:"type:varchar(200);not null;default:'';"`
}

type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Brands   Brands
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`
	URL   string `gorm:"type:varchar(200);not null"`
	Index int32  `gorm:"type:int;not null;default:1"`
}

type Goods struct {
	BaseModel

	// 外键
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"not null;default:false;"`
	ShipFree bool `gorm:"not null;default:false;comment:'免运费'"`
	IsNew    bool `gorm:"not null;default:false;comment:'新品'"`
	IsHot    bool `gorm:"not null;default:false;comment:'热卖'"`

	Name    string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`

	ClickNum int32 `gorm:"type:int;not null;default:0"`
	SoldNum  int32 `gorm:"type:int;not null;default:0"`
	FavNum   int32 `gorm:"type:int;not null;default:0"`

	// 不是 int 类型
	MarketPrice float32 `gorm:"not null"`
	ShopPrice   float32 `gorm:"not null"`
	GoodsBrief  string  `gorm:"type:varchar(100);not null;comment:'商品简介'"`

	Images          GormList `gorm:"type:varchar(1000);not null"`
	DescImages      GormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`
}
