package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"go-api-server/api"
	. "go-api-server/common"
	"go-api-server/config"
)

func Handler_Test(resp http.ResponseWriter, req *http.Request) {
	Log("Request: ", req.URL)

	resp.Write([]byte("Ok"))
}

var file_server_handler = http.FileServer(http.Dir(config.API_server_html_dir))

func Handler_HTML(resp http.ResponseWriter, req *http.Request) {
	Log("Request: ", req.URL)

	resp.Header().Set("Cache-Control", config.HTTP_header_no_cache)
	file_server_handler.ServeHTTP(resp, req)
}

func Handler_APIJS(resp http.ResponseWriter, req *http.Request) {
	Log("Request: ", req.URL)

	resp.Header().Set("Access-Control-Allow-Origin", config.CORS_header_origin)
	resp.Header().Set("Content-Type", "text/javascript")
	resp.Header().Set("Cache-Control", config.HTTP_header_no_cache)
	resp.Write(api.JS_Source.Bytes())
}

func Handler_Call(resp http.ResponseWriter, req *http.Request) {

	if req.Method == "OPTIONS" {
		resp.Header().Set("Access-Control-Allow-Origin", config.CORS_header_origin)
		resp.Header().Set("Access-Control-Allow-Headers", config.CORS_header_headers)
		resp.WriteHeader(http.StatusOK)
		return
	}

	Log("Request:", req.URL)
	defer Log_Time(time.Now())

	resp.Header().Set("Access-Control-Allow-Origin", config.CORS_header_origin)
	resp.Header().Set("Cache-Control", config.HTTP_header_no_cache)
	resp.Header().Set("Content-Type", "application/json")

	var err error
	defer func() {
		if err != nil {
			Log("Error:", err.Error())
			message := strings.ReplaceAll(err.Error(), "\"", "")
			message = fmt.Sprintf(`{"Error": "%s"}`, message)
			http.Error(resp, message, http.StatusInternalServerError)
		}
	}()

	funcname := req.PathValue("funcname")
	api_handler, ok := api.Call_Handler_Map[funcname]
	if !ok {
		err = fmt.Errorf("no such method in api: %s", funcname)
		return
	}

	api_argument := reflect.New(api_handler.Argument)

	req.Body = http.MaxBytesReader(resp, req.Body, api.Max_Request_Size)
	json_decoder := json.NewDecoder(req.Body)
	json_decoder.DisallowUnknownFields()

	if err = json_decoder.Decode(api_argument.Interface()); err != nil {
		return
	}

	call_return := api_handler.Function.Call([]reflect.Value{reflect.ValueOf(req.Context()), api_argument})
	if !call_return[1].IsNil() {
		err = call_return[1].Interface().(error)
		return
	}

	if call_return[0].Elem().Kind() == reflect.String {
		result := call_return[0].Elem().Interface().(string)
		Log(result)
		resp.Write([]byte(result))
		return
	}

	if call_return[0].Elem().Type() == reflect.TypeOf([]byte(nil)) {
		result := call_return[0].Elem().Interface().([]byte)
		Log(string(result))
		resp.Write(result)
		return
	}

	var encode_buffer bytes.Buffer
	json_encoder := json.NewEncoder(&encode_buffer)
	json_encoder.SetEscapeHTML(false)
	if err = json_encoder.Encode(call_return[0].Interface()); err != nil {
		return
	}

	Log(encode_buffer.String())
	resp.Write(encode_buffer.Bytes())
}
