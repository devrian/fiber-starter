package responses

import (
	"fiber-starter/app/model"
	"time"
)

type UserResponse struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Code        string    `json:"code"`
	Role        string    `json:"role"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Img         string    `json:"img"`
	Status      bool      `json:"status"`
	OTP         string    `json:"otp"`
	Version     int       `json:"version"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

func (r *UserResponse) Transform(data model.User, otpToken string) {
	r.Name = data.Name
	r.Email = data.Email
	r.Code = data.Code
	r.Role = data.Role
	r.Phone = data.Phone
	r.Address = data.Address.String
	r.Img = data.Img.String
	r.Status = data.Status
	r.OTP = otpToken
	r.Version = int(data.Version)
	r.CreatedDate = data.CreatedDate
	r.UpdatedDate = data.UpdatedDate.Time
}
