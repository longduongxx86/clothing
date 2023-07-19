package incoming

type AccountLogin struct {
	Email    string `json:"email" form:"email" param:"email" `
	Password string `json:"password" form:"password" param:"password"`
}

type AccountResetPassword struct {
	Email       string `json:"email" form:"email" param:"email" `
	NewPassword string `json:"new-password" form:"new-password" param:"new-password"`
}

type AccountSignIn struct {
	Fullname    string `json:"fullname" form:"fullname" param:"fullname" `
	Email       string `json:"email" form:"email" param:"email" `
	PhoneNumber string `json:"phone_number" form:"phone_number" param:"phone_number" `
	Address     string `json:"address" form:"address" param:"address" `
	Password    string `json:"password" form:"password" param:"password" `
}
