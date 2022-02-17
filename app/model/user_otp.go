package model

import (
	"database/sql"
	"fiber-starter/pkg/common"
	"fiber-starter/pkg/utils"
	"math/rand"
	"time"
)

const (
	TOKEN_EXPIRED_TIME = 3 // In minutes
	OTP_VIA_EMAIL      = "email"
	OTP_VIA_PHONE      = "phone"
)

type UserOTP struct {
	ID          int64        `db:"id"`
	Email       string       `db:"email"`
	Phone       string       `db:"phone"`
	OTP         string       `db:"otp"`
	ExpiredDate time.Time    `db:"expired_date"`
	CreatedDate time.Time    `db:"created_date"`
	UpdatedDate sql.NullTime `db:"updated_date"`
}

func (u UserOTP) GenerateOTP() (string, error) {
	rand.Seed(time.Now().UnixNano())
	otp, err := common.GenerateString(`[\d]{6}`)

	return otp, err
}

func (u UserOTP) GetExpiredDate() time.Time {
	now := time.Now().In(time.UTC)
	return now.Add(time.Minute * TOKEN_EXPIRED_TIME)
}

func ViaValidMailPhoneChannel(channel, emailPhone string) (string, bool) {
	var via string
	if common.IsEmail(emailPhone) {
		via = OTP_VIA_EMAIL
	}

	_, _, validPhone := common.IsPhone(emailPhone)
	if validPhone {
		via = OTP_VIA_PHONE
	}

	if len(via) <= 0 {
		return via, false
	}

	if channel == utils.ChannelApp {
		if via == OTP_VIA_EMAIL && common.IsEmail(emailPhone) {
			return via, true
		}
		if via == OTP_VIA_PHONE {
			_, _, validPhone := common.IsPhone(emailPhone)
			if validPhone {
				return via, true
			}
		}
	} else if channel == utils.ChannelWeb {
		if via == OTP_VIA_EMAIL && common.IsEmail(emailPhone) {
			return via, true
		}
	}

	return via, false
}
