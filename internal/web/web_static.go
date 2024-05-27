//go:build static

package web

import (
	"embed"
	"net/http"
	"strings"

	"github.com/ItsNotGoodName/ipcmanview/pkg/chiext"
)

//go:generate pnpm install
//go:generate pnpm run build

//go:embed all:dist
var fs embed.FS

func FS(skipPaths ...string) func(next http.Handler) http.Handler {
	return chiext.StaticEmbedFS(chiext.StaticFSConfig{
		FileSystem: fs,
		Root:       "dist",
		SPA:        true,
		Redirect: func(r *http.Request) bool {
			for _, path := range skipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					return false
				}
			}
			return true
		},
	})
}
