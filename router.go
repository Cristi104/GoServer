package main

import (
	"net/http"
	"regexp"
	"strings"
)

var routes = []route{
	makeRoute("GET", "/", MainHandler),
	makeRoute("GET", "/signin", SignInPostHandler),
	makeRoute("POST", "/signin", SignInGetHandler),
	makeRoute("GET", "/signup", SignUpPostHandler),
	makeRoute("POST", "/signup", SignUpGetHandler),
	makeRoute("GET", "/home", HomeHandler),
	makeRoute("GET", "/data/conversations", DataConversationsGetHandler),
	makeRoute("POST", "/data/messages/load", DataMessagesLoadHandler),
	makeRoute("POST", "/data/messages/add", DataMessagesAddHandler),
	makeRoute("GET", "/data/account", DataAccountHandler),
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func makeRoute(method string, pattern string, handler http.HandlerFunc) route {
	return route{
		method:  method,
		regex:   regexp.MustCompile("^" + pattern + "$"),
		handler: handler,
	}
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			route.handler(w, r)
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}
