package handlers_test

import (
    "encoding/json"
    "fmt"
    "github.com/google/go-cmp/cmp"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/streadway/amqp"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func TestGetAllWorkflows(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
    var method string
    var path string
    var request *http.Request
    var response *httptest.ResponseRecorder
    var expectedCode int
    var workflows []models.Workflow
    var isFound bool
    var workflowAux models.Workflow

    workflow = models.Workflow{}

    // A json object of the workflow input.
    data = `{"array":[1],"boolean":true,"float":1.1,"integer":1,"object":{"key":"value"},"string":"string"}`

    err = json.Unmarshal([]byte(data), &workflow.Data)

    if err != nil {
        t.Fatalf("Failed to parse the JSON-encoded data of the workflow data %s: %s", data, err.Error())
    }

    workflow.Steps = []string{utils.GenerateRandomString(10)}

    body = fmt.Sprintf(`{"data":%s,"steps":["%s"]}`, data, workflow.Steps[0])

    workflow, err = s.Datastore.CreateWorkflow(workflow)

    if err != nil {
        t.Fatalf("Failed to create a new workflow with %s: %s", body, err.Error())
    }

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
    }

    t.Logf("Workflow: %s", string(bodyBytes))

    method = "GET"

    path = "/workflows"

    request, err = http.NewRequest(method, path, nil)

    if err != nil {
        t.Fatalf("Failed to create the request: %s", err.Error())
    }

    t.Logf("Request: method=%s and path=%s", method, path)

    response = httptest.NewRecorder()

    r.ServeHTTP(response, request)

    expectedCode = http.StatusOK

    if expectedCode != response.Code {
        t.Errorf("Test failed, the expected response code: %d, got: %d", expectedCode, response.Code)
        return
    }

    err = json.NewDecoder(response.Body).Decode(&workflows)

    if err != nil {
        t.Fatalf("Failed to parse the JSON response body: %s", err.Error())
    }

    isFound = false

    for _, workflowAux = range workflows {
        // Evaluate the equality of the simulated data with those returned from the associated functionality.
        if cmp.Equal(workflow, workflowAux) {
            isFound = true
            break
        }
    }

    if !isFound {
        t.Errorf("Test failed, the workflow not found in the list of all workflows: %s", string(bodyBytes))
        return
    }

    t.Logf("Test successful, the workflow found in the list of all workflows: %s", string(bodyBytes))
}

func TestConsumeWorkflow(t *testing.T) {
    var timeout <-chan time.Time
    var isTimedOut bool
    var body string
    var err error
    var workflow models.Workflow
    var data string
    var bodyBytes []byte
    var message amqp.Publishing
    var method string
    var path string
    var request *http.Request
    var response *httptest.ResponseRecorder
    var expectedCode int
    var workflowAux models.Workflow
    var bodyBytesAux []byte

    // Consume.
    timeout = time.After(60 * time.Second)

    // Keep trying until we're timed out or consumed all workflows.
    for {
        select {
        // Got a timeout!
        case <-timeout:
            isTimedOut = true
            break
        default:
        }

        if isTimedOut {
            break
        }

        body, err = s.MessageBroker.Consume(s.MessageBroker.Queue.Name, true)

        if err != nil {
            t.Fatalf("Failed to consume a workflow from the queue: %s", err.Error())
        }

        // In case of the returned body is empty it means that there isn't any workflow to be consumed anymore.
        if body == "" {
            break
        }
    }

    workflow = models.Workflow{}

    // A json object of the workflow input.
    data = `{"array":[1],"boolean":true,"float":1.1,"integer":1,"object":{"key":"value"},"string":"string"}`

    err = json.Unmarshal([]byte(data), &workflow.Data)

    if err != nil {
        t.Fatalf("Failed to parse the JSON-encoded data of the workflow data %s: %s", data, err.Error())
    }

    workflow.Steps = []string{utils.GenerateRandomString(10)}

    body = fmt.Sprintf(`{"data":%s,"steps":["%s"]}`, data, workflow.Steps[0])

    workflow, err = s.Datastore.CreateWorkflow(workflow)

    if err != nil {
        t.Fatalf("Failed to create a new workflow with %s: %s", body, err.Error())
    }

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
    }

    t.Logf("Workflow: %s", string(bodyBytes))

    // Publish the workflow.
    message = amqp.Publishing{
        Body: bodyBytes,
    }

    err = s.MessageBroker.Publish("events", "random-key", false, false, message)

    if err != nil {
        t.Fatalf("Failed to publish the workflow %+v on the queue: %s", workflow, err.Error())
    }

    // Consume the workflow.
    method = "GET"

    path = "/workflows/consume"

    request, err = http.NewRequest(method, path, nil)

    if err != nil {
        t.Fatalf("Failed to create the request: %s", err.Error())
    }

    request.Header.Set("Content-Type", "application/json")

    t.Logf("Request: method=%s and path=%s", method, path)

    response = httptest.NewRecorder()

    r.ServeHTTP(response, request)

    expectedCode = http.StatusOK

    if expectedCode != response.Code {
        t.Errorf("Test failed, response: code=%d and body=%+v", response.Code, response.Body)
        return
    }

    err = json.NewDecoder(response.Body).Decode(&workflowAux)

    if err != nil {
        t.Fatalf("Failed to parse the JSON response body: %s", err.Error())
    }

    // Evaluate the equality of the simulated data with those returned from the associated functionality.
    if !cmp.Equal(workflow, workflowAux) {
        bodyBytesAux, err = json.Marshal(workflowAux)

        if err != nil {
            t.Fatalf("Failed to obtain the JSON encoding of the returned workflow %+v: %s", workflowAux, err.Error())
        }

        t.Errorf("Test failed, the expected workflow returned: %s, got: %s", string(bodyBytes), string(bodyBytesAux))
        return
    }

    t.Logf("Test successful, response: code=%d and body=%s", response.Code, string(bodyBytes))
}
