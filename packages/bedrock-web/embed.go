package renderer

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var FS embed.FS

func GetHttpFS() http.FileSystem {
	fs, _ := fs.Sub(FS, "dist")

	return http.FS(fs)
}
