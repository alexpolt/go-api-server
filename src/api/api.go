package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"sort"
	"strings"
	"text/template"

	. "go-api-server/common"
	"go-api-server/config"
	"go-api-server/lang"
)

const Max_Request_Size = 64 * 1024

var api_handlers_list = map[string]reflect.Value{

	"load_projects":  reflect.ValueOf(Load_Projects),
	"get_project":    reflect.ValueOf(Get_Project),
	"create_project": reflect.ValueOf(Create_Project),

	"load_namespaces":  reflect.ValueOf(Load_Namespaces),
	"get_namespace":    reflect.ValueOf(Get_Namespace),
	"create_namespace": reflect.ValueOf(Create_Namespace),
	"delete_namespace": reflect.ValueOf(Delete_Namespace),

}

type Call_Handler struct {
	Function reflect.Value
	Argument reflect.Type
}

var JS_Source bytes.Buffer
var Call_Handler_Map = make(map[string]Call_Handler)

func init() {
	for handler_name, value := range api_handlers_list {
		handler_type := value.Type()

		if !(handler_type.NumIn() == 2 &&
			handler_type.In(0) == reflect.TypeOf((*context.Context)(nil)).Elem() &&
			handler_type.In(1).Kind() == reflect.Pointer) {
			Panic(errors.New("wrong call handler arguments"))
		}
		if !(handler_type.NumOut() == 2 &&
			handler_type.Out(0).Kind() == reflect.Pointer &&
			handler_type.Out(1) == reflect.TypeOf((*error)(nil)).Elem()) {
			Panic(errors.New("wrong call handler return values"))
		}

		handler := Call_Handler{
			Function: value,
			Argument: handler_type.In(1).Elem(),
		}

		Call_Handler_Map[handler_name] = handler
	}

	Make_JS_Source()
}

var js_template = `
	globalThis.api = globalThis.API = {}

	api.logging = false
	api.client_api_path = "{{.client_api_path}}"
	api.server_api_path = "{{.server_api_path}}"
	api.api_path = typeof window === "undefined" ? api.server_api_path : api.client_api_path
	
	api.make_call = async (name, arg) => {
		if(API.logging) console.log("api", name, arg)
		const opts = {
			method: "POST",
			cache: "no-cache",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(arg),
		}
		const reply = await fetch(api.api_path + name, opts)
    const data = await reply.json()
    if(!reply.ok) {
      if(typeof data !== "object") {
        throw new Error(` + "`[${name}]: ${reply.status}: ${reply.statusText}`" + `)
      }
      if(data.Error) {
        throw new Error(` + "`[${name}] (${reply.status}): ${data.Error}`" + `)
      }
    }
    return data
	}

	api._check = (n, arg, orig) => {
		if(typeof orig !== "object" && typeof orig !== typeof arg) {
			throw new Error("[" + n + "] Expected different type: " + (typeof arg) + " vs " + (typeof orig))
		}
		if(orig instanceof Array && !(arg instanceof Array)) {
			throw new Error("[" + n + "] Expected array: " + (typeof arg))
		}
		if(arg instanceof Array) {
			return
		}
		for(const key of Object.keys(orig)) {
			if(! (key in arg))
				throw new Error("[" + n + "] No required [" + key + "] in argument", orig)
		}
		for(const key of Object.keys(arg)) {
			// if(key.startsWith("__")) continue
			// if(typeof arg[key] === "function") continue
			if(! (key in orig))
				throw new Error("[" + n + "] Invalid [" + key + "] in argument", orig)
			if(typeof arg[key] !== typeof orig[key])
				throw new Error("[" + n + "] Types of [" + key + "] in argument differ: " + (typeof arg[key]) + " vs " + (typeof orig[key]))
			if(arg[key] instanceof Object) {
					if(key === "Content" && arg[key].Type) {
						if(!api.make_content_node[arg[key].Type])
							throw new Error("[" + n + "] Invalid [" + key + "].Type (" + arg[key].Type + ") in argument")
						const content = api.make_content_node[arg[key].Type]()
						api._check(n, arg[key], content)
					} else {
						api._check(n, arg[key], orig[key])
					}
			}
		}
	}

	api._arg = {
		{{range $n, $v := .arguments}}
		"{{$n}}": {{$v}},{{end}}
	}

	api.fn = {
		{{range $n, $v := .arguments}}
		"{{$n}}": async (arg) => { api._check("{{$n}}", arg, api._arg["{{$n}}"]); return api.make_call("{{$n}}", arg) },{{end}}
	}

	api.make_arg = {
		{{range $n, $v := .arguments}}
		"{{$n}}": () => { return JSON.parse(` + "`{{$v}}`" + `) },{{end}}
	}

`

func Make_JS_Source() {
	tmpl, err := template.New("").Parse(js_template)
	Panic(err)

	arguments := make(map[string]string)
	var argument_json []byte
	for name, handler := range Call_Handler_Map {
		if handler.Argument.Kind() == reflect.Slice {
			argument_json, err = json.Marshal(reflect.MakeSlice(handler.Argument, 0, 0).Interface())
		} else {
			argument_json, err = json.Marshal(reflect.New(handler.Argument).Elem().Interface())
		}
		Panic(err)
		arguments[name] = string(argument_json)
	}


	data := map[string]any{
		"client_api_path":   config.API_client_url,
		"server_api_path":   config.API_server_url,
		"arguments":         arguments,
	}

	tmpl.Execute(io.Writer(&JS_Source), data)
}
