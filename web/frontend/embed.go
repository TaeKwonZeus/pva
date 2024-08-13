package frontend

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var frontendEmbed embed.FS

var Embed = func() fs.FS {
	f, err := fs.Sub(frontendEmbed, "dist")
	if err != nil {
		panic(err)
	}
	return f
}()
