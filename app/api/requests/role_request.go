package requests

type (
	RoleCreateRequest struct {
		Name   string `json:"name" validate:"required"`
		Status *bool  `json:"status" validate:"required"`
	}

	RoleUpdateRequest struct {
		Name    string `json:"name" validate:"required"`
		Status  *bool  `json:"status" validate:"required"`
		Version int    `json:"version" validate:"required"`
	}
)
