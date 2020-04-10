package routes

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
)

type Route struct {
    Name string
    Method string
    Pattern string
    HandlerFunc http.HandlerFunc
}

// Configure the routes of the API endpoints.
func ConfigureRoutes(r *mux.Router, s *server.Server) *mux.Router {
    var routeList []Route

    // It refers to the operation linked to the service status.
    routeList = append(routeList, AddGetStatusRoute())
 
    // It refers to the operations linked to workflow(s).
    routeList = append(routeList, AddGetAllWorkflowsRoute(s))
    routeList = append(routeList, AddCreateWorkflowRoute(s))
    routeList = append(routeList, AddUpdateWorkflowRoute(s))
    routeList = append(routeList, AddConsumeWorkflowRoute(s))

    for _, route := range routeList {
        r.Name(route.Name).
            Methods(route.Method).
            Path(route.Pattern).
            HandlerFunc(route.HandlerFunc)
    }

    return r
}
