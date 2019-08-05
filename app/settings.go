package app

import (
	"github.com/jmoiron/sqlx"
)

const (
	VERSION = "0.1.0"
)

var (
	DB   *sqlx.DB
	Conf map[string]string
)
