package util

import "net/url"

func Get(values url.Values, name, def string) string {
	if values.Has(name) {
		return values.Get(name)
	}
	return def
}
