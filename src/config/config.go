package config

import (
	"os"

	. "go-api-server/common"
)

var (
	APP_path = Default_string(os.Getenv("APP_PATH"), "/app")

	API_server_port    string = Default_string(os.Getenv("API_SERVER_PORT"), "8080")
	API_server_address string = Default_string(os.Getenv("API_SERVER_HOST"), "") + ":" + API_server_port
	API_client_address string = "localhost:" + API_server_port

	API_server_url string = "http://" + API_server_address + "/call/"
	API_client_url string = "http://" + API_client_address + "/call/"

	API_server_html_dir string = APP_path + "/src" // FileServer adds "/html/" from URL

	HTTP_header_no_cache = "no-store"
	CORS_header_origin   = "*"
	CORS_header_headers  = "Content-Type, Accept"
)
