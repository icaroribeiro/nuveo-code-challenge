package postgresdb_test

import (
    "encoding/json"
    "fmt"
    "github.com/google/go-cmp/cmp"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "testing"
)

func TestGetAllWorkflows(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
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

    workflow, err = datastore.CreateWorkflow(workflow)

    if err != nil {
        t.Fatalf("Failed to create a new workflow with %s: %s", body, err.Error())
    }

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
    }

    t.Logf("Workflow: %s", string(bodyBytes))

    workflows, err = datastore.GetAllWorkflows()

    if err != nil {
        t.Fatalf("Failed to get the list of all workflows: %s", err.Error())
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

func TestGetWorkflow(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
    var workflowAux models.Workflow
    var bodyBytesAux []byte

    workflow = models.Workflow{}

    // A json object of the workflow input.
    data = `{"array":[1],"boolean":true,"float":1.1,"integer":1,"object":{"key":"value"},"string":"string"}`

    err = json.Unmarshal([]byte(data), &workflow.Data)

    if err != nil {
        t.Fatalf("Failed to parse the JSON-encoded data of the workflow data %s: %s", data, err.Error())
    }

    workflow.Steps = []string{utils.GenerateRandomString(10)}

    body = fmt.Sprintf(`{"data":%s,"steps":["%s"]}`, data, workflow.Steps[0])

    workflow, err = datastore.CreateWorkflow(workflow)

    if err != nil {
        t.Fatalf("Failed to create a new workflow with %s: %s", body, err.Error())
    }

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
    }

    t.Logf("Workflow: %s", string(bodyBytes))

    workflowAux, err = datastore.GetWorkflow(workflow.ID)

    if err != nil {
        t.Fatalf("Failed to get the workflow with the id %s: %s", workflow.ID, err.Error())
    }

    // Evaluate the equality of the simulated data with those returned from the associated functionality.
    if !(cmp.Equal(workflow, workflowAux)) {
        bodyBytesAux, err = json.Marshal(workflowAux)

        if err != nil {
            t.Fatalf("Failed to obtain the JSON encoding of the returned workflow %+v: %s", workflowAux, err.Error())
        }

        t.Errorf("Test failed, the expected workflow returned: %s, got: %s", string(bodyBytes), string(bodyBytesAux))
        return
    }

    t.Logf("Test successful, the returned workflow: %s", string(bodyBytes))
}
