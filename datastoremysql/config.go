package datastoremysql

import "time"

type Config struct {
	DSNString       string
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
