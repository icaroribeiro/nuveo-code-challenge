package postgresdb

import (
    "database/sql"
    "database/sql/driver"
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/lib/pq"
)

const (
    QueryGetAllWorkflows = "GetAllWorkflows"
    QueryGetWorkflow = "GetWorkflow"
)

func AddGetAllWorkflowsStatement(unpreparedStmts map[string]string) {
    unpreparedStmts[QueryGetAllWorkflows] = `
            SELECT id, status, data, steps
            FROM workflows;
        `
}

func AddGetWorkflowStatement(unpreparedStmts map[string]string) {
    unpreparedStmts[QueryGetWorkflow] = `
            SELECT id, status, data, steps
            FROM workflows
            WHERE id = $1;
        `
}

func (d *Datastore) GetAllWorkflows() ([]models.Workflow, error) {
    var rows *sql.Rows
    var err error
    var workflows []models.Workflow
    var workflow models.Workflow
    var dataArray []sql.NullString
    var data sql.NullString
    var name string   
    var driverValue driver.Value
    var isOK bool

    rows, err = d.Stmts[QueryGetAllWorkflows].Query()

    if err != nil {
        return workflows, err
    }

    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&workflow.ID, 
                    &workflow.Status, 
                    &workflow.Data, 
                    pq.Array(&dataArray))

        // Complete the list of the names of all workflow steps.
        workflow.Steps = []string{}

        for _, data = range dataArray {
            driverValue, err = data.Value()
    
            name, isOK = driverValue.(string)
    
            if !isOK {
                return workflows, 
                    fmt.Errorf("it wasn't possible to get data from the list of the names of all workflow steps")
            }
    
            workflow.Steps = append(workflow.Steps, name)
        }

        if err != nil {
            return workflows, err
        }

        workflows = append(workflows, workflow)
    }

    return workflows, nil
}

func (d *Datastore) GetWorkflow(id string) (models.Workflow, error) {
    var err error
    var workflow models.Workflow
    var dataArray []sql.NullString
    var data sql.NullString
    var name string   
    var driverValue driver.Value
    var isOK bool

    err = d.Stmts[QueryGetWorkflow].QueryRow(id).Scan(&workflow.ID, 
                                                &workflow.Status, 
                                                &workflow.Data, 
                                                pq.Array(&dataArray))
    
    // Complete the list of the names of all workflow steps.
    for _, data = range dataArray {
        driverValue, err = data.Value()

        name, isOK = driverValue.(string)

        if !isOK {
            return workflow, fmt.Errorf("it wasn't possible to get data from the list of the names of all workflow steps")
        }

        workflow.Steps = append(workflow.Steps, name)
    }
    
    if err != nil {
        if err != sql.ErrNoRows {
            return workflow, err
        }
    }

    return workflow, nil
}
