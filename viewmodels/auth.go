package viewmodels

//SignIn used for binding sign in form values
type SignIn struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

//SignUp used for binding sign up form values
type SignUp struct {
	Email                string `form:"email" json:"email" binding:"required"`
	Password             string `form:"password" json:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" binding:"required"`
}
