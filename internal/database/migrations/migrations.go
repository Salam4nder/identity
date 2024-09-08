package migrations

import (
	"embed"
	_ "embed"
)

//go:embed *.sql
var Files embed.FS
