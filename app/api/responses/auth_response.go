package responses

import (
	"fiber-starter/app/model"
)

type RegisterResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Code  string `json:"code"`
	Role  string `json:"role"`
	OTP   string `json:"otp"`
}

func (r *RegisterResponse) Transform(data model.User, otpToken string) {
	r.Name = data.Name
	r.Email = data.Email
	r.Phone = data.Phone
	r.Code = data.Code
	r.Role = data.Role
	r.OTP = otpToken
}

type LoginResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Code   string `json:"code"`
	Token  string `json:"token"`
	Role   string `json:"role"`
	Img    string `json:"img"`
	Status bool   `json:"status"`
	OTP    string `json:"otp"`
}

func (r *LoginResponse) Transform(data model.User, otpToken string) {
	r.Name = data.Name
	r.Email = data.Email
	r.Phone = data.Phone
	r.Code = data.Code
	r.Token = data.RememberToken.String
	r.Role = data.Role
	r.Img = data.Img.String
	r.Status = data.Status
	r.OTP = otpToken
}

type SendOTPTokenResponse struct {
	OTPToken   string `json:"otp_token"`
	Status     bool   `json:"status"`
	EmailPhone string `json:"email_phone"`
}

func (r *SendOTPTokenResponse) Transform(otpToken string, status bool, emailPhone string) {
	r.OTPToken = otpToken
	r.Status = status
	r.EmailPhone = emailPhone
}
