package translator

import "net/url"
// javascript "encodeURI()"
// so we embed js to our golang programm
func encodeURI(s string) string {
	return url.QueryEscape(s)
}