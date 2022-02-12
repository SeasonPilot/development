package forms

type PassWordLoginForm struct {
	Name      string `form:"name" json:"name" binding:"required,mobile"`
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5" `
	CaptchaID string `form:"captcha_id" json:"captcha_id" binding:"required"`
}
