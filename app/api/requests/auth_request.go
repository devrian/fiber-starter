package requests

type (
	RegisterRequest struct {
		Name            string `json:"name" validate:"required"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
		Phone           string `json:"phone" validate:"required"`
	}

	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	ResetRequest struct {
		EmailPhone string `json:"email_phone" validate:"required"`
	}

	ResetPasswordRequest struct {
		EmailPhoneToken string `json:"email_phone_token" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required"`
	}

	ValidateOTPTokenRequest struct {
		EmailPhone string `json:"email_phone" validate:"required"`
		OTPToken   string `json:"otp_token" validate:"required"`
	}
)
