package request

type User struct {
	Username string `json:"username" binding:"omitempty,min=6,max=20"`
	Password string `json:"password" binding:"omitempty,min=6,max=64"`
	Nickname string `json:"nickname" binding:"omitempty,min=6,max=20"`
	Avatar   string `json:"avatar" binding:"omitempty,url"`
	Phone    string `json:"phone" binding:"omitempty,number,min=11,max=11"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
