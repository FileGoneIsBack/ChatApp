package server

import "login/core/models/functions"

var TemplateCache *functions.TemplateCache

type Config struct {
	Addr      string
	Secure    bool
	Cert, Key string
}
