package forms

type PassWordLoginForm struct {
	Name     string `form:"name" json:"name" binding:"required,mobile"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
}
