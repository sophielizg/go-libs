package datastoremysql

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Connection struct {
	Config Config
	db     *sqlx.DB
}

func (c *Connection) Open() error {
	db, err := sqlx.Connect("mysql", c.Config.DSNString)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.db = db
	return nil
}

func (c *Connection) Close() {
	c.db.Close()
}

// Expose underlying db to query directly
func (c *Connection) Db() *sqlx.DB {
	return c.db
}
