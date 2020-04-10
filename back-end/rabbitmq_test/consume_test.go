package rabbitmq_test

import (    
    "encoding/json"
    "fmt"
    "github.com/google/go-cmp/cmp"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/streadway/amqp"
    "testing"
    "time"
)

func TestConsume(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
    var confirms chan amqp.Confirmation
    var message amqp.Publishing
    var confirmed amqp.Confirmation
    var workflowAux models.Workflow
    var bodyBytesAux []byte
    var timeout <-chan time.Time
    var isTimedOut bool

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

    err = messageQueue.Chan.Confirm(false)

	if err != nil {
        t.Fatalf("Failed to configure the channel in confirm mode: %s", err.Error())
	}

    confirms = messageQueue.Chan.NotifyPublish(make(chan amqp.Confirmation, 1))
        
    message = amqp.Publishing{
        Body: bodyBytes,
    }

    err = messageQueue.Publish("events", "random-key", false, false, message)

    if err != nil {
        t.Fatalf("Failed to publish the workflow %+v on the queue: %s", workflow, err.Error())
    }

    confirmed = <-confirms

    if !confirmed.Ack {
        t.Errorf("Test failed, the publisher couldn't confirm the delivery of the delivery tag: %d", confirmed.DeliveryTag)
        return
    }

    t.Logf("Confirmed the delivery with the delivery tag: %d", confirmed.DeliveryTag)

    // Consume.
    timeout = time.After(60 * time.Second)
    
    // Keep trying until we're timed out, got an error or consumed the related workflow.
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

        body, err = messageQueue.Consume(messageQueue.Queue.Name, true)

        if err != nil {
            t.Fatalf("Failed to consume a workflow from the queue: %s", err.Error())
        }

        err = json.Unmarshal([]byte(body), &workflowAux)

        if err != nil {
            t.Fatalf("Failed to parse the JSON-encoded data of the workflow %s from the queue: %s", body, err.Error())
        }

        if (workflow.ID == workflowAux.ID) {
            break
        }
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

    t.Logf("Test successful, the consumed workflow: %s", string(bodyBytes))
}
