package routes

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/handlers"
)

func AddGetStatusRoute() Route {
    var route = Route {
            Name: "GetStatus",
            Method: "GET",
            Pattern: "/status",
            HandlerFunc: handlers.GetStatus(),
        }

    return route
}
