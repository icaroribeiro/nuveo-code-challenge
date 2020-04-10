package rabbitmq

import (
    "github.com/streadway/amqp"
)

func (mb *MessageBroker) Publish(exchange, key string, mandatory, immediate bool, message amqp.Publishing) error {
    var err error

    // Send a message from the client to an exchange on the server.
    err = mb.Chan.Publish(exchange, key, mandatory, immediate, message)

    if err != nil {
        return err
    }

    return nil
}
