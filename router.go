package main

import (
	"fmt"
	"net/http"
	"strings"
)

//Context : ...
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Params map[string]string
}

//HandlerFunc : alternative of http.HandlerFunc
type HandlerFunc func(*Context)

//Router : alternative of http's Handler
type Router struct {
	hMaps map[string]map[string]HandlerFunc
}

//Handle : alternative of http.HandleFunc
func (router *Router) Handle(method, pattern string, handler HandlerFunc) {
	if router.hMaps == nil {
		router.hMaps = make(map[string]map[string]HandlerFunc)
	}
	hMap, ext := router.hMaps[method]
	if !ext {
		hMap = make(map[string]HandlerFunc)
		router.hMaps[method] = hMap
	}

	hMap[pattern] = handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hMap, ext := router.hMaps[r.Method]

	if !ext {
		http.NotFound(w, r)
		return
	}

	for pattern, handler := range hMap {
		if params, matching := match(strings.Split(r.URL.Path, "/"), r.URL.Path, pattern); matching {
			context := Context{Writer: w, Request: r, Params: params}
			handler(&context)
			if e := recover(); e != nil {
				fmt.Println(e)
			}
			r.Body.Close()
			break
		}
	}
}

func match(path []string, rawpath, pattern string) (map[string]string, bool) {
	if pattern == rawpath {
		return nil, true
	}

	patterns := strings.Split(pattern, "/")
	ofs := 0
	if strings.Contains(pattern, "~") {
		if len(patterns) > len(path)+1 {
			return nil, false
		}
		if len(patterns)-len(path) == 1 {
			ofs = 1
		}
	} else if len(patterns) != len(path) {
		return nil, false
	}

	params := make(map[string]string)
	for i := 0; i < len(patterns)-ofs; i++ {
		switch {
		case patterns[i] == path[i]:
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			params[patterns[i][1:]] = path[i]
		case len(patterns[i]) > 0 && patterns[i][0] == '~':
			params[patterns[i][1:]] = strings.Join(path[i:], "/")
			break
		default:
			return nil, false
		}
	}

	return params, true
}
