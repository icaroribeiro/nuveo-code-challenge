package rabbitmq

import (
    "fmt"
    "github.com/streadway/amqp"
)

func (mb *MessageBroker) Consume(queue string, autoAck bool) (string, error) {
    var message amqp.Delivery
    var isOK bool
    var err error

    // Receive synchronously a single Delivery from the head of the queue from the server.
    message, isOK, err = mb.Chan.Get(queue, autoAck)

    if err != nil {
        return "", err
    }

    if !isOK {
        return "", fmt.Errorf("there is no delivery waiting on the queue")
    }

    return string(message.Body), nil
}
