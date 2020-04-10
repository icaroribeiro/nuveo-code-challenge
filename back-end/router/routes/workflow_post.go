package routes

import (    
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/handlers"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/middlewares"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
)

func AddCreateWorkflowRoute(s *server.Server) Route {
    var route = Route {
            Name: "CreateWorkflow",
            Method: "POST",
            Pattern: "/workflow",
            HandlerFunc: middlewares.AdaptFunc(handlers.CreateWorkflow(s)).
                With(middlewares.ValidateRequestHeaderFields(map[string]string{
                        "Content-Type": "application/json",
                    }),
            ),
        }

    return route
}
