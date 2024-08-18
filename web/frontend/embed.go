package frontend

import (
	"embed"
	"io/fs"
	"sync"
)

//go:embed dist
var dist embed.FS

var getEmbed = sync.OnceValue(func() fs.FS {
	f, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return f
})

func Embed() fs.FS {
	return getEmbed()
}
