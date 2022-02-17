package repository

import (
	"fiber-starter/app/model"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gosimple/slug"
)

type RoleRepository interface {
	Insert(dbctx db.DBCtx, rl model.Role) (model.Role, error)
	Update(dbctx db.DBCtx, rl model.Role) error
	Delete(dbctx db.DBCtx, code, deletedBy string) error
	GetAll(dbctx db.DBCtx, f model.RoleFilter) ([]model.Role, error)
	GetAllTotal(dbctx db.DBCtx, f model.RoleFilter) (int64, error)
	GetByCode(dbctx db.DBCtx, code string) (model.Role, error)
	GetIDBySlug(dbctx db.DBCtx, slug string) (int64, error)
	GetVersionByCode(dbctx db.DBCtx, code string) (int32, error)
}

type roleRepository struct {
}

func NewRoleRepository() *roleRepository {
	return &roleRepository{}
}

func (r *roleRepository) Insert(dbctx db.DBCtx, rl model.Role) (model.Role, error) {
	var ID int64

	rl.Version = 1
	rl.CreatedDate = time.Now().In(time.UTC)
	rl.Slug = slug.Make(string(rl.Name))

	paramQ := []interface{}{rl.Code, rl.Name, rl.Slug, rl.Status, rl.CreatedDate, rl.CreatedBy, rl.Version}
	q := `insert into roles (code, name, slug, status, created_date, created_by, version) values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := dbctx.TX.QueryRow(dbctx.Ctx, q, paramQ...).Scan(&ID)
	rl.ID = ID

	return rl, err
}

func (r *roleRepository) Update(dbctx db.DBCtx, rl model.Role) error {
	paramQ := []interface{}{rl.Name, rl.Slug, rl.Status, rl.UpdatedDate.Time, rl.UpdatedBy.String, rl.Version, rl.Code}
	q := `update roles set name = $1, slug = $2, status = $3, updated_date = $4, updated_by = $5, version = $6 where deleted_date is null and code = $7`
	_, err := dbctx.TX.Exec(dbctx.Ctx, q, paramQ...)

	return err
}

func (r *roleRepository) Delete(dbctx db.DBCtx, code, deletedBy string) error {
	timeStamp := time.Now().In(time.UTC)
	_, err := dbctx.TX.Exec(dbctx.Ctx,
		"update roles set updated_date = $1, updated_by = $2, deleted_date = $3, deleted_by = $4, status = $5 where code = $6",
		timeStamp, deletedBy, timeStamp, deletedBy, false, code,
	)

	return err
}

func (r *roleRepository) GetAll(dbctx db.DBCtx, f model.RoleFilter) ([]model.Role, error) {
	var roles []model.Role

	q := `select * from roles where deleted_date is null`
	q += f.WhereGetAll()
	q += fmt.Sprintf(` %s`, common.OrderByOffsetLimitSQL(f.Paging))

	err := pgxscan.Select(dbctx.Ctx, dbctx.DB, &roles, q)

	return roles, err
}

func (r *roleRepository) GetAllTotal(dbctx db.DBCtx, f model.RoleFilter) (int64, error) {
	var total int64

	q := `select count(*) from roles where deleted_date is null`
	q += f.WhereGetAll()

	err := pgxscan.Get(dbctx.Ctx, dbctx.DB, &total, q)

	return total, err
}

func (r *roleRepository) GetByCode(dbctx db.DBCtx, code string) (model.Role, error) {
	var rl model.Role

	q := `select id, code, name, slug, status, created_date, created_by, updated_date, updated_by, version from roles where deleted_date is null and code = $1 limit 1`
	err := dbctx.DB.QueryRow(dbctx.Ctx, q, code).Scan(&rl.ID, &rl.Code, &rl.Name, &rl.Slug, &rl.Status, &rl.CreatedDate, &rl.CreatedBy, &rl.UpdatedDate, &rl.UpdatedBy, &rl.Version)

	return rl, err
}

func (r *roleRepository) GetIDBySlug(dbctx db.DBCtx, slug string) (int64, error) {
	var ID int64

	q := fmt.Sprintf(`select id from roles where deleted_date is null and status = %v and slug = $1 limit 1`, true)
	err := dbctx.DB.QueryRow(dbctx.Ctx, q, slug).Scan(&ID)

	return ID, err
}

func (r *roleRepository) GetVersionByCode(dbctx db.DBCtx, code string) (int32, error) {
	var v int32

	q := `select version from roles where deleted_date is null and code = $1 limit 1`
	err := dbctx.DB.QueryRow(dbctx.Ctx, q, code).Scan(&v)

	return v, err
}
