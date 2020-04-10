package handlers

import (
    "encoding/json"
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/streadway/amqp"
    "net/http"
)

func CreateWorkflow(s *server.Server) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var err error
        var workflow models.Workflow
        var bodyBytes []byte
        var body string
        var i int
        var step string
        var message amqp.Publishing

        err = json.NewDecoder(r.Body).Decode(&workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to decode the request body: %s", err.Error())})
            return
        }

        if workflow.Data.IsEmpty() == true {
            utils.RespondWithJson(w, http.StatusBadRequest, 
                map[string]string{"error": "The data field is required and must be set to a non-empty json object"})
            return
        }

        if len(workflow.Steps) == 0 {
            utils.RespondWithJson(w, http.StatusBadRequest, 
                map[string]string{"error": "The steps field is required and must be set to an array containing " + 
                    "at least one name of workflow step"})
            return
        }

        bodyBytes, err = json.Marshal(workflow.Data)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to obtain the JSON encoding of the workflow data %+v: %s", 
                    workflow.Data, err.Error())})
            return
        }

        body = fmt.Sprintf(`{"data":%s`, string(bodyBytes))

        for i, step = range workflow.Steps {
            if i == 0 {
                body += fmt.Sprintf(`,"steps":["%s"`, step)
            } else {
                body += fmt.Sprintf(`,"%s"`, step)
            }
        }

        body += `]}`

        workflow, err = s.Datastore.CreateWorkflow(workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to create a new workflow with %s: %s", body, err.Error())})
            return
        }

        bodyBytes, err = json.Marshal(workflow)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to obtain the JSON encoding of the workflow %+v: %s", 
                    workflow, err.Error())})
            return
        }

        // Create a message to be sent to the queue.
        message = amqp.Publishing{
            Body: bodyBytes,
        }

        err = s.MessageBroker.Publish("events", "random-key", false, false, message)

        if err != nil {
            utils.RespondWithJson(w, http.StatusInternalServerError, 
                map[string]string{"error": fmt.Sprintf("Failed to publish the workflow %+v on the queue: %s", 
                    workflow, err.Error())})
            return
        }

        utils.RespondWithJson(w, http.StatusCreated, workflow)
    })
}
