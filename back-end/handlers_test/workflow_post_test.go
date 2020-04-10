package handlers_test

import (
    "encoding/json"
    "fmt"
    "github.com/google/go-cmp/cmp"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestCreateWorkflow(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error   
    var method string
    var path string
    var body string
    var request *http.Request
    var response *httptest.ResponseRecorder
    var expectedCode int
    var workflowAux models.Workflow
    var bodyBytes []byte
    var bodyBytesAux []byte

    workflow = models.Workflow{}

    // A json object of the workflow input.
    data = `{"array":[1],"boolean":true,"float":1.1,"integer":1,"object":{"key":"value"},"string":"string"}`

    err = json.Unmarshal([]byte(data), &workflow.Data)

    if err != nil {
        t.Fatalf("Failed to parse the JSON-encoded data of the workflow data %s: %s", data, err.Error())
    }

    workflow.Steps = []string{utils.GenerateRandomString(10)}

    method = "POST"

    path = "/workflow"

    body = fmt.Sprintf(`{"data":%s,"steps":["%s"]}`, data, workflow.Steps[0])

    request, err = http.NewRequest(method, path, strings.NewReader(body))

    if err != nil {
        t.Fatalf("Failed to create the request: %s", err.Error())
    }

    request.Header.Set("Content-Type", "application/json")

    t.Logf("Request: method=%s, path=%s and body=%s", method, path, body)

    response = httptest.NewRecorder()

    r.ServeHTTP(response, request)

    expectedCode = http.StatusCreated

    if expectedCode != response.Code {
        t.Errorf("Test failed, response: code=%d and body=%+v", response.Code, response.Body)
        return
    }

    err = json.NewDecoder(response.Body).Decode(&workflowAux)

    if err != nil {
        t.Fatalf("Failed to parse the JSON response body: %s", err.Error())
    }

    workflow.ID = workflowAux.ID
    workflow.Status = workflowAux.Status

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
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
