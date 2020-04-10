package rabbitmq

import (
    "fmt"
    "github.com/streadway/amqp"
)

// The MBConfig stores all the parameters to configure the message broker settings.
type MBConfig struct {  
    Username string
    Password string
    Host     string
    Port     string
    Name     string
}

// The MessageBroker groups all the variables necessary to connect and work with messages 
// by means of a channel used to interface with our backing queue.
type MessageBroker struct {
    Conn *amqp.Connection
    Chan *amqp.Channel
    Queue amqp.Queue
}

func InitializeMB(mbConfig MBConfig) (MessageBroker, error) {
    var connString string
    var conn *amqp.Connection
    var err error
    var channel *amqp.Channel
    var queue amqp.Queue

    // Set up the connection string of the message queue.
    connString = SetUpConnString(mbConfig)

    // Connect to the RabbitMQ instance.
    conn, err = amqp.Dial(connString)

    if err != nil {
        return MessageBroker{}, fmt.Errorf("it wasn't possible to establish connection with RabbitMQ: %s", err.Error())
    }

    // Create a channel from the connection to access the data in tue queue.
    channel, err = conn.Channel()

    if err != nil {
        return MessageBroker{}, fmt.Errorf("it wasn't possible to open RabbitMQ channel: %s", err.Error())
    }

    // Declare an exchange on the server thath will bind to the queue to send and receive messages.
    err = channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)

    if err != nil {
        return MessageBroker{}, err
    }

    // Declare a queue to hold messages and deliver to consumers.
    queue, err = channel.QueueDeclare(mbConfig.Name, true, false, false, false, nil)

    if err != nil {
        return MessageBroker{}, fmt.Errorf("it wasn't possible to declare que queue: %s", err.Error())
    }

    // Bind the queue to the exchange to send and receive data from the queue.
    err = channel.QueueBind(mbConfig.Name, "#", "events", false, nil)

    if err != nil {
        return MessageBroker{}, fmt.Errorf("it wasn't possible to bind to the queue: %s", err.Error())
    }

    return MessageBroker{Conn: conn, Chan: channel, Queue: queue}, err
}

// It builds the connection string of the message queue.
func SetUpConnString(mbConfig MBConfig) string {
    var connString string

    connString = fmt.Sprintf("amqp://%s:%s@%s:%s",
                        mbConfig.Username,
                        mbConfig.Password,
                        mbConfig.Host,
                        mbConfig.Port)

    return connString
}

// It closes the message queue, releasing any open resources.
func (mq *MessageBroker) Close() error {
    var err error

    err = mq.Chan.Close()

    if err != nil {
        return err
    }

    err = mq.Conn.Close()

    if err != nil {
        return err
    }

    return nil
}
