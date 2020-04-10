package rabbitmq_test

import (
    "encoding/json"
    "fmt"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/models"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/streadway/amqp"
    "testing"
)

func TestPublish(t *testing.T) {
    var workflow models.Workflow
    var data string
    var err error
    var body string
    var bodyBytes []byte
    var confirms chan amqp.Confirmation
    var message amqp.Publishing
    var confirmed amqp.Confirmation

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

    // Enable pusblishing confirms by putting the channel into confirm mode 
    // so that to ensure all publishings have successfully been received by the server. 
    err = messageQueue.Chan.Confirm(false)

	if err != nil {
        t.Fatalf("Failed to configure the channel in confirm mode: %s", err.Error())
	}

    // Registers a listener for reliable publishing.
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

    t.Logf("Test successful, the publisher confirmed the delivery with the delivery tag: %d", confirmed.DeliveryTag)
}
