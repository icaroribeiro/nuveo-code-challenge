package postgresdb_test

import (
    "encoding/json"
    "fmt"
    "github.com/google/go-cmp/cmp"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "testing"
)

func TestCreateWorkflow(t *testing.T) {
    var workflow models.Workflow
    var data string
    var body string
    var workflowAux models.Workflow
    var err error
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

    body = fmt.Sprintf(`{"data":%s,"steps":["%s"]}`, data, workflow.Steps[0])

    t.Logf("Workflow: %s", body)

    workflowAux, err = datastore.CreateWorkflow(workflow)

    if err != nil {
        t.Fatalf("Failed to create a new workflow with %s: %s", body, err.Error())
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

    t.Logf("Test successful, the created workflow: %s", string(bodyBytes))
}
