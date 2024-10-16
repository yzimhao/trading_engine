package migrations

import "embed"

//go:embed postgres/*.sql
var Migrations embed.FS
