package chiext

import (
	"io/fs"
	"net/http"
	"strings"
)

// func init() {
// 	mime.AddExtensionType(".js", "application/javascript")
// }

type StaticFSConfig struct {
	FileSystem fs.FS
	Root       string
	SPA        bool
	Redirect   func(r *http.Request) bool
}

// StaticEmbedFS adds GET handlers for all files and folders using the given filesystem.
func StaticEmbedFS(config StaticFSConfig) func(next http.Handler) http.Handler {
	if config.Redirect == nil {
		config.Redirect = func(r *http.Request) bool { return false }
	}
	if config.Root != "" {
		fsys, err := fs.Sub(config.FileSystem, config.Root)
		if err != nil {
			panic(err)
		}
		config.FileSystem = fsys
	}

	fsHandler := http.StripPrefix("/", http.FileServer(http.FS(config.FileSystem)))
	indexHandler := func(w http.ResponseWriter, r *http.Request) {
		index, err := http.FS(config.FileSystem).Open("/index.html")
		if err != nil {
			panic(err)
		}
		defer index.Close()

		stat, err := index.Stat()
		if err != nil {
			panic(err)
		}

		http.ServeContent(w, r, "index.html", stat.ModTime(), index)
	}

	files, err := fs.ReadDir(config.FileSystem, ".")
	if err != nil {
		panic(err)
	}

	routes := []string{}
	for _, f := range files {
		routes = append(routes, "/"+f.Name())
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				indexHandler(w, r)
				return
			}
			for _, route := range routes {
				if strings.HasPrefix(r.URL.Path, route) {
					if r.URL.Path == "/index.html" {
						indexHandler(w, r)
						return
					}

					fsHandler.ServeHTTP(w, r)
					return
				}
			}

			if config.Redirect(r) {
				r.URL.Path = "/"
				indexHandler(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
