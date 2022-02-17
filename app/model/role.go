package model

import (
	"database/sql"
	"fiber-starter/pkg/common"
	"fiber-starter/pkg/utils"
	"fmt"
	"strings"
	"time"
)

const (
	ROLE_ADMIN  = "admin"
	ROLE_CUST   = "customer"
	ROLE_PREFIX = "role"
)

var roleSlugList = []string{ROLE_ADMIN, ROLE_CUST}

type Role struct {
	ID          int64          `db:"id"`
	Code        string         `db:"code"`
	Name        string         `db:"name"`
	Slug        string         `db:"slug"`
	Status      bool           `db:"status"`
	CreatedDate time.Time      `db:"created_date"`
	CreatedBy   string         `db:"created_by"`
	UpdatedDate sql.NullTime   `db:"updated_date"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	DeletedDate sql.NullTime   `db:"deleted_date"`
	DeletedBy   sql.NullString `db:"deleted_by"`
	Version     int32          `db:"version"`
}

type RoleFilter struct {
	Status string
	Search string
	Paging common.PaginateQueryOffset
}

func (f RoleFilter) WhereGetAll() string {
	var q string

	if len(f.Search) > 0 {
		var orWhere []string
		orWhere = append(orWhere, fmt.Sprintf(`lower(code) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(name) like lower('%%%s%%')`, f.Search))
		orWhere = append(orWhere, fmt.Sprintf(`lower(slug) like lower('%%%s%%')`, f.Search))

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

func CheckValidChannelRole(channel string, role string) error {
	if channel == utils.ChannelWeb && role != ROLE_ADMIN {
		return fmt.Errorf("%s", `channel role is invalid`)
	}

	return nil
}

func CheckValidRoleSlug(slug string) error {
	var err error
	if !common.CheckStringContains(slug, roleSlugList) {
		statusStr := strings.Join(roleSlugList, ", ")
		errMsg := fmt.Sprintf(`invalid role arguments in action not includes in %s.`, statusStr)
		err = fmt.Errorf("%s", errMsg)
	}

	return err
}
