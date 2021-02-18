package embed

import (
	"embed"
)

//go:embed web
// StaticFiles ... files that will be embedded in the application
var StaticFiles embed.FS
