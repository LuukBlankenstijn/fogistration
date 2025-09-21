//go:build !docker_prod

package spa

import (
	"io/fs"
)

func SpaFs() (fs.FS, error) {
	return nil, nil
}
