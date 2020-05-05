package plugin

import "net/http"

type SkipperFunc func(r *http.Request) bool

// Default: check result is false, skip this step.
var DefaultSkipperFunc = func(r *http.Request) bool {
	return false
}
