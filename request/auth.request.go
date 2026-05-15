package request

type LoginRequest struct {
	Email    string `json:"email" form:"required, email" binding:"required,email"`
	Password string `json:"password" form:"required" binding:"required"`
}

type RegisterRequest struct {
	Name            string `json:"name" form:"required" binding:"required"`
	Email           string `json:"email" form:"required, email" binding:"required,email"`
	Password        string `json:"password" form:"required,min=6" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" form:"required,eqfield=Password" binding:"required,eqfield=Password"`
}
