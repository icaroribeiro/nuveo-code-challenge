package handlers

import (
    "encoding/json"
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/services"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "net/http"
)

func GetAllWorkflows(s *server.Server) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var err error        
        var workflows []models.Workflow

        workflows, err = s.Datastore.GetAllWorkflows()

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to get the list of all workflows: %s", err.Error())})
            return
        }

        utils.RespondWithJson(w, http.StatusOK, workflows)
    })
}

func ConsumeWorkflow(s *server.Server) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var body string
        var err error
        var workflow models.Workflow
        var jsonBytes []byte

        body, err = s.MessageBroker.Consume(s.MessageBroker.Queue.Name, true)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to consume a workflow from the queue: %s", err.Error())})
            return
        }

        if body == "" {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to consume a workflow from the queue: " + 
                "there is no delivery waiting on the queue")})
            return
        }

        err = json.Unmarshal([]byte(body), &workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to parse the JSON-encoded data of " +
                    "the workflow %s from the queue: %s", body, err.Error())})
            return
        }

        jsonBytes, err = json.Marshal(workflow.Data)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to obtain the JSON encoding of the workflow.Data %+v: %s",
                    workflow.Data, err.Error())})
            return
        }

        // Save json related to workflow.Data in a csv file.
        err = services.GenerateCSVFile(jsonBytes, s.StorageDir, workflow.ID)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to generate a csv file with the data of " + 
                    "the workflow %s from the queue: %s", body, err.Error())})
            return
        }

        utils.RespondWithJson(w, http.StatusOK, workflow)
    })
}
