package embed

import (
	"embed"
)

//go:embed web
// StaticFiles ... files that will be embedded with the application
var StaticFiles embed.FS
