package routes

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/handlers"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
)

func AddGetAllWorkflowsRoute(s *server.Server) Route {
    var route = Route {
        Name: "GetAllWorkflows",
        Method: "GET",
        Pattern: "/workflows",
        HandlerFunc: handlers.GetAllWorkflows(s),
    }

    return route
}

func AddConsumeWorkflowRoute(s *server.Server) Route {
    var route = Route {
        Name: "ConsumeWorkflow",
        Method: "GET",
        Pattern: "/workflows/consume",
        HandlerFunc: handlers.ConsumeWorkflow(s),
    }

    return route
}
