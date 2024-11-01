package main

import (
	"net/http"
)

var HTTP_Routes = map[string]http.HandlerFunc{

	"/": Handler_Test,

	"GET /html/": Handler_HTML,

	"GET /api.js": Handler_APIJS,

	"/call/{funcname}": Handler_Call,
}
