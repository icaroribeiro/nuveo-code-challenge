package routes

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/handlers"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/middlewares"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
)

func AddUpdateWorkflowRoute(s *server.Server) Route {
    var route = Route {
            Name: "UpdateWorkflow",
            Method: "PATCH",
            Pattern: "/workflows/{workflowId}",
            HandlerFunc: middlewares.AdaptFunc(handlers.UpdateWorkflow(s)).
                With(middlewares.ValidateRequestHeaderFields(map[string]string{
                        "Content-Type": "application/json",
                    }),
            ),
        }

    return route
}
