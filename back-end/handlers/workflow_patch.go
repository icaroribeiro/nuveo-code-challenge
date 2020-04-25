package handlers

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "net/http"
)

func UpdateWorkflow(s *server.Server) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var params map[string]string
        var workflowId string
        var err error
        var workflow models.Workflow
        var body string    
        var nRowsAffected int64

        params = mux.Vars(r)

        workflowId = params["workflowId"]

        if workflowId == "" {
            utils.RespondWithJson(w, http.StatusBadRequest, 
                map[string]string{"error": "The id is required and must be set to a non-empty value in the request URL"})
            return
        }

        workflow, err = s.Datastore.GetWorkflow(workflowId)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to get the workflow with the id %s: %s", workflowId, err.Error())})
            return
        }

        err = json.NewDecoder(r.Body).Decode(&workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to decode the request body: %s", err.Error())})
            return
        }

        workflow.ID = workflowId

        if workflow.Status != "consumed" {
            utils.RespondWithJson(w, http.StatusBadRequest, 
                map[string]string{"error": `The status field is required and must be set to "consumed"`})
            return
        }

        body = fmt.Sprintf(`{"status":"%s"}`, workflow.Status)

        nRowsAffected, err = s.Datastore.UpdateWorkflow(workflowId, workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to update the workflow with the id %s with %s: %s", 
                    workflowId, body, err.Error())})
            return
        }

        if nRowsAffected == 0 {
            utils.RespondWithJson(w, http.StatusConflict, 
                map[string]string{"error": fmt.Sprintf("Failed to update the workflow with the id %s with %s: " + 
                    "the workflow wasn't found", workflowId, body)})
            return
        }

        if nRowsAffected != 1 {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to update the workflow with the id %s with %s: the expected number of"+
                "workflows updated: %d, got: %d", workflowId, body, 1, nRowsAffected)})
            return
        }

        utils.RespondWithJson(w, http.StatusOK, workflow)
    })
}
