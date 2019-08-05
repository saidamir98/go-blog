package app

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

const (
	VERSION = "0.1.0"
)

var (
	Mux  *mux.Router
	DB   *sqlx.DB
	Conf map[string]string
)
