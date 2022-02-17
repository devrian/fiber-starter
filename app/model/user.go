package model

import (
	"database/sql"
	"fiber-starter/pkg/common"
	"fmt"
	"strings"
	"time"
)

const (
	USER_PREFIX  = "user"
	DEFAULT_PASS = "12345678"
)

type User struct {
	ID            int64          `db:"id"`
	Code          string         `db:"code"`
	RoleID        int32          `db:"role_id"`
	Role          string         `db:"role"`
	Name          string         `db:"name"`
	Email         string         `db:"email"`
	Phone         string         `db:"phone"`
	Password      string         `db:"password"`
	Address       sql.NullString `db:"address"`
	Img           sql.NullString `db:"img"`
	RememberToken sql.NullString `db:"remember_token"`
	Status        bool           `db:"status"`
	CreatedDate   time.Time      `db:"created_date"`
	CreatedBy     string         `db:"created_by"`
	UpdatedDate   sql.NullTime   `db:"updated_date"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	DeletedDate   sql.NullTime   `db:"deleted_date"`
	DeletedBy     sql.NullString `db:"deleted_by"`
	Version       int32          `db:"version"`
}

type UserFilter struct {
	Roles  []string
	Status string
	Search string
	Paging common.PaginateQueryOffset
}

func (f UserFilter) WhereGetAll() string {
	var q string

	if len(f.Roles) > 0 {
		var roles []string
		for _, role := range f.Roles {
			roles = append(roles, common.AddQuote(role))
		}
		q += fmt.Sprintf(` and role in (%s)`, strings.Join(roles, ", "))
	}

	if len(f.Search) > 0 {
		var orWhere []string
		orWhere = append(orWhere, fmt.Sprintf(`lower(code) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(role) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(name) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(email) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(phone) like lower('%%%s%%')`, f.Search))

		q += fmt.Sprintf(` and (%s) `, strings.Join(orWhere, ` or `))
	}

	if len(f.Status) > 0 {
		var status bool
		if f.Status == "true" {
			status = true
		}
		q += fmt.Sprintf(` and status = %v `, status)
	}

	return q
}
