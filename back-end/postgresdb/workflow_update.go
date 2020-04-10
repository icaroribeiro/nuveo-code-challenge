package postgresdb

import (
    "database/sql"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
)

const (
    QueryUpdateWorkflow = "UpdateWorkflow"
)

func AddUpdateWorkflowStatement(unpreparedStmts map[string]string) {
    unpreparedStmts[QueryUpdateWorkflow] = `
            UPDATE workflows
            SET
            status = $1
            WHERE id = $2;
        `
}

func (d *Datastore) UpdateWorkflow(id string, workflow *models.Workflow) (int64, error) {
    var result sql.Result
    var err error
    var nRowsAffected int64

    result, err = d.Stmts[QueryUpdateWorkflow].Exec(workflow.Status, id)

    if err != nil {
        return 0, err
    }

    nRowsAffected, err = result.RowsAffected()

    if err != nil {
        return 0, err
    }

    (*workflow).ID = id

    return nRowsAffected, nil
}
