package db

import (
	"context"
	"fiber-starter/config"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBCtx struct {
	Ctx context.Context
	DB  *pgxpool.Conn
	TX  pgx.Tx
}

func (d *DBCtx) Set(ctx context.Context, db *pgxpool.Conn, tx pgx.Tx) {
	d.Ctx = ctx
	d.DB = db
	d.TX = tx
}

// Init
func Init(c *config.Config) *pgxpool.Pool {
	db := c.Database
	pass := "postgres"
	if len(db.Password) > 0 {
		pass = db.Password
	}

	dsn := fmt.Sprintf(`%s://%s:%s@%s:%s/%s`, db.Driver, db.User, pass, db.Host, db.Port, db.Name)
	conn, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}

	return conn
}

// Parsing Error
func ParseErr(err error) string {
	switch pqe := err.(type) {
	case *pgconn.PgError:

		// error duplicate
		if pqe.Code == "23505" {
			m := strings.ReplaceAll(pqe.Detail, "(", "")
			m = strings.ReplaceAll(m, ")", "")
			m = strings.ReplaceAll(m, "=", " ")
			m = strings.ReplaceAll(m, "Key", "")
			m = strings.ReplaceAll(m, "=", "")
			return m
		}
	}
	return err.Error()
}
