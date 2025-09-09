package helper

import (
	"github.com/jackc/pgx/v5"
)

func IsNoRowsError(err error) bool {
	return err == pgx.ErrNoRows
}
