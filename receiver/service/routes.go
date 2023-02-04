package service

import (
	"github.com/julienschmidt/httprouter"
)

// route :
type route struct {
	Method string            //HTTP method
	Path   string            //url endpoint
	Handle httprouter.Handle //Controller function which dispatches the right HTML page and/or data for each route
}

var routes = []route{
	route{
		"POST",
		"/json",
		receiveJSON,
	},
	route{
		"GET",
		"/ready",
		ready,
	},
}
