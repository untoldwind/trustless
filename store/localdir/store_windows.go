package localdir

import (
	"net/url"
	"path/filepath"
)

func pathFromURL(url *url.URL) string {
	path := url.Path

	if len(path) > 2 && path[0] == '/' && path[2] == ':' {
		return filepath.FromSlash(path[1:])
	}
	return filepath.FromSlash(path)
}
