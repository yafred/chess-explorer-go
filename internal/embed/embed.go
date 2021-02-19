package embed

import (
	"embed"
)

//go:embed css
//go:embed img
//go:embed javascript
//go:embed index.html
// StaticFiles ... files that will be embedded in the application
var StaticFiles embed.FS
