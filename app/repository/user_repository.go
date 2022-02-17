package repository

import (
	"fiber-starter/app/model"
	"fiber-starter/db"
	"fiber-starter/pkg/common"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type UserRepository interface {
	Insert(dbctx db.DBCtx, u model.User) (model.User, error)
	Update(dbctx db.DBCtx, u model.User) error
	Delete(dbctx db.DBCtx, code, deletedBy string) error
	UpdatePasswordByEmailOrPhone(dbctx db.DBCtx, password, emailPhone string) error
	UpdateStatusByEmailOrPhone(dbctx db.DBCtx, status bool, emailPhone string) error
	GetByCode(dbctx db.DBCtx, code string) (model.User, error)
	GetByEmail(dbctx db.DBCtx, email string) (model.User, error)
	GetByEmailOrPhone(dbctx db.DBCtx, emailPhone string) (model.User, error)
	GetAll(dbctx db.DBCtx, f model.UserFilter) ([]model.User, error)
	GetAllTotal(dbctx db.DBCtx, f model.UserFilter) (int64, error)
	GetVersionByCode(dbctx db.DBCtx, code string) (int32, error)
}

type userRepository struct {
}

func NewUserRepository() *userRepository {
	return &userRepository{}
}

func (r *userRepository) Insert(dbctx db.DBCtx, u model.User) (model.User, error) {
	var ID int64
	u.CreatedDate = time.Now().In(time.UTC)

	paramQ := []interface{}{u.Code, u.RoleID, u.Role, u.Name, u.Email, u.Phone, u.Password, u.Address.String, u.Img.String, u.RememberToken.String, u.Status, u.CreatedDate, u.CreatedBy, 1}

	q := `insert into users (code, role_id, role, name, email, phone, password, address, img, remember_token, status, created_date, created_by, version) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) returning id`
	err := dbctx.TX.QueryRow(dbctx.Ctx, q, paramQ...).Scan(&ID)
	if err != nil {
		return u, err
	}
	u.ID = ID

	return u, err
}

func (r *userRepository) Update(dbctx db.DBCtx, u model.User) error {
	paramQ := []interface{}{u.Name, u.Phone, u.Address.String, u.Img.String, u.Status, u.UpdatedDate.Time, u.UpdatedBy.String, u.Version, u.Code}
	q := `update users set name = $1, phone = $2, address = $3, img = $4, status = $5, updated_date = $6, updated_by = $7, version = $8 where deleted_date is null and code = $9`
	_, err := dbctx.TX.Exec(dbctx.Ctx, q, paramQ...)

	return err
}

func (r *userRepository) Delete(dbctx db.DBCtx, code, deletedBy string) error {
	timeStamp := time.Now().In(time.UTC)
	_, err := dbctx.TX.Exec(dbctx.Ctx,
		"update users set updated_date = $1, updated_by = $2, deleted_date = $3, deleted_by = $4, status = $5 where code = $6",
		timeStamp, deletedBy, timeStamp, deletedBy, false, code,
	)

	return err
}

func (r *userRepository) UpdatePasswordByEmailOrPhone(dbctx db.DBCtx, password, emailPhone string) error {
	paramQ := []interface{}{password, time.Now().In(time.UTC), emailPhone, emailPhone}

	q := `update users set password = $1, updated_date = $2 where deleted_date is null and (email = $3 or phone = $4)`
	exec, err := dbctx.TX.Exec(dbctx.Ctx, q, paramQ...)
	if exec.RowsAffected() <= 0 {
		return fmt.Errorf(`%s`, "update password failed")
	}

	return err
}

func (r *userRepository) UpdateStatusByEmailOrPhone(dbctx db.DBCtx, status bool, emailPhone string) error {
	paramQ := []interface{}{status, nil, time.Now().In(time.UTC), emailPhone, emailPhone}

	q := `update users set status = $1, remember_token = $2, updated_date = $3 where deleted_date is null and (email = $4 or phone = $5)`
	exec, err := dbctx.TX.Exec(dbctx.Ctx, q, paramQ...)
	if exec.RowsAffected() <= 0 {
		return fmt.Errorf(`%s`, "update status failed")
	}

	return err
}

func (r *userRepository) GetByCode(dbctx db.DBCtx, code string) (model.User, error) {
	var u model.User

	q := `select * from users where deleted_date is null and code = $1 limit 1`
	err := pgxscan.Get(dbctx.Ctx, dbctx.DB, &u, q, code)

	return u, err
}

func (r *userRepository) GetByEmail(dbctx db.DBCtx, email string) (model.User, error) {
	var u model.User

	q := `select * from users where deleted_date is null and email = $1 limit 1`
	err := pgxscan.Get(dbctx.Ctx, dbctx.DB, &u, q, email)

	return u, err
}

func (r *userRepository) GetByEmailOrPhone(dbctx db.DBCtx, emailPhone string) (model.User, error) {
	var u model.User

	q := `select * from users where deleted_date is null and (email = $1 or phone = $2) limit 1`
	err := pgxscan.Get(dbctx.Ctx, dbctx.DB, &u, q, emailPhone, emailPhone)

	return u, err
}

func (r *userRepository) GetAll(dbctx db.DBCtx, f model.UserFilter) ([]model.User, error) {
	var users []model.User

	q := `select * from users where deleted_date is null`
	q += f.WhereGetAll()
	q += fmt.Sprintf(` %s`, common.OrderByOffsetLimitSQL(f.Paging))

	err := pgxscan.Select(dbctx.Ctx, dbctx.DB, &users, q)

	return users, err
}

func (r *userRepository) GetAllTotal(dbctx db.DBCtx, f model.UserFilter) (int64, error) {
	var total int64

	q := `select count(*) from users where deleted_date is null`
	q += f.WhereGetAll()

	err := pgxscan.Get(dbctx.Ctx, dbctx.DB, &total, q)

	return total, err
}

func (r *userRepository) GetVersionByCode(dbctx db.DBCtx, code string) (int32, error) {
	var v int32

	q := `select version from users where deleted_date is null and code = $1 limit 1`
	err := dbctx.DB.QueryRow(dbctx.Ctx, q, code).Scan(&v)

	return v, err
}
