package forms

type CreateOrderForm struct {
	Name   string `json:"name" binding:"required"`
	Mobile string `json:"mobile" binding:"required,mobile"`
	Post   string `json:"post" binding:"required"`
	Addr   string `json:"addr" binding:"required"`
}
