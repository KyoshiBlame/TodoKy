package core_http_server

import "net/http"


type Route struct {
	Method string
	Path string
	Handler http.HandlerFunc
}

//Механика которая даёт возможность каждой фиче определять свой набор роутов, методой, путей и хендлеров
func NewRoute(
	Method string,
	Path string,
	Handler http.HandlerFunc,
) *Route {
	return &Route {
		Method: Method,
		Path: Path,
		Handler: Handler,
	}
}