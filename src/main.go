package main

import (
	"net/http"
	"strings"

	"go-api-server/api"
	. "go-api-server/common"
	"go-api-server/config"
	"go-api-server/db"
)

const ()

var ()

func main() {
	Log("Staring server on ", config.API_server_address)
	Log("File server dir ", config.API_server_html_dir)

	var name_list []string
	for name := range api.Call_Handler_Map {
		name_list = append(name_list, name)
	}
	Log("Registered api handlers: ", strings.Join(name_list, "; "))

	name_list = nil
	for route, handler := range HTTP_Routes {
		name_list = append(name_list, route)
		http.HandleFunc(route, handler)
	}
	Log("Registered routes: ", strings.Join(name_list, "; "))

	db.Connect()
	defer db.Close()

	http.ListenAndServe(config.API_server_address, nil)
}
