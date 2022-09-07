package testdata

import (
	"context"
	"database/sql"
)

// Object is the object.
// go-gen-opt:gen
type Object struct {
	ctx      context.Context
	required string  `go-gen-opt:"required"`
	optional int     `go-gen-opt-default:"5"`
	conn     *sql.DB `go-gen-opt-func:"createConnection" go-gen-opt:"required"`
}

func createConnection(host string, port int) (*sql.DB, error) {
	c := Config{
		Required: "asdf",
		Host:     "asdf",
		Port:     324,
	}
	_ = c
	return nil, nil
}
