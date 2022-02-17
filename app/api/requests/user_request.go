package requests

type (
	UserCreateRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Phone    string `json:"phone" validate:"required"`
		RoleCode string `json:"role_code"`
		Address  string `json:"address"`
		Img      string `json:"img"`
	}

	UserUpdateRequest struct {
		IsProfile *bool  `json:"is_profile" validate:"required"`
		Name      string `json:"name" validate:"required"`
		Phone     string `json:"phone" validate:"required"`
		Address   string `json:"address"`
		Img       string `json:"img"`
		Status    *bool  `json:"status"`
		Version   int    `json:"version" validate:"required"`
	}
)
