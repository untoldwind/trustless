package localdir

import "net/url"

func pathFromURL(url *url.URL) string {
	return url.Path
}
