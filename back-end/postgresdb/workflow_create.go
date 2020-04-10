package postgresdb

import (
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/lib/pq"
)

const (
    QueryCreateWorkflow = "CreateWorkflow"
)

func AddCreateWorkflowStatement(unpreparedStmts map[string]string) {
    unpreparedStmts[QueryCreateWorkflow] = `
            INSERT INTO 
            workflows (data, steps)
            VALUES ($1, $2) RETURNING id, status;
        `
}

func (d *Datastore) CreateWorkflow(workflow models.Workflow) (models.Workflow, error) {
    var err error

    err = d.Stmts[QueryCreateWorkflow].QueryRow(workflow.Data, pq.Array(workflow.Steps)).Scan(&workflow.ID, &workflow.Status)

    if err != nil {
        return workflow, err
    }

    if workflow.ID == "" {
        return workflow, fmt.Errorf("it wasn't possible to get the id of the generated record")
    }

    return workflow, nil
}
