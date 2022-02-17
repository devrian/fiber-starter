package responses

import (
	"fiber-starter/app/model"
)

type RoleResponse struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Slug        string `json:"slug"`
	Status      bool   `json:"status"`
	Version     int    `json:"version"`
	CreatedDate string `json:"created_date"`
	CreatedBy   string `json:"created_by"`
	UpdatedDate string `json:"updated_date"`
	UpdatedBy   string `json:"updated_by"`
}

func (r *RoleResponse) Transform(data model.Role) {
	r.Name = data.Name
	r.Code = data.Code
	r.Slug = data.Slug
	r.Status = data.Status
	r.Version = int(data.Version)
	r.CreatedBy = data.CreatedBy
	r.UpdatedBy = data.UpdatedBy.String
	if !data.CreatedDate.IsZero() {
		r.CreatedDate = data.CreatedDate.Format("2006-01-02 15:04:05")
	}
	if data.UpdatedDate.Valid {
		r.UpdatedDate = data.UpdatedDate.Time.Format("2006-01-02 15:04:05")
	}
}
