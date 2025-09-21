//go:build docker_prod

package spa

import (
	"embed"
	"io/fs"
)

//go:embed all:frontend/dist
var staticFs embed.FS

func SpaFs() (fs.FS, error) {
	return fs.Sub(staticFs, "frontend/dist")
}
