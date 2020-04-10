package postgresdb_test

import (
    "encoding/json"
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "testing"
)

func TestUpdateWorkflow(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
    var nRowsAffected int64

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

    workflow.Status = "consumed"

    body = fmt.Sprintf(`{"status":%s}`, workflow.Status)

    t.Logf("New workflow data: %s", body)

    nRowsAffected, err = datastore.UpdateWorkflow(workflow.ID, &workflow)

    if err != nil {
        t.Fatalf("Failed to update the workflow with the id %s with %s: %s", workflow.ID, body, err.Error())
    }

    if nRowsAffected == 0 {
        t.Errorf("Test failed, the workflow with the id %s wasn't found", workflow.ID)
    }

    if nRowsAffected != 1 {
        t.Errorf("Test failed, the expected number of workflows updated: %d, got: %d", 1, nRowsAffected)
        return
    }

    bodyBytes, err = json.Marshal(workflow)

    if err != nil {
        t.Fatalf("Failed to obtain the JSON encoding of the workflow %+v: %s", workflow, err.Error())
    }

    t.Logf("Test successful, the updated workflow: %s", string(bodyBytes))
}
