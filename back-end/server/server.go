package server

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
)

// This structure is an abstraction of the server that allows to "attach" some resources in order to make them
// available during the API requests. Here, it's used to store other structure that holds attributes to manage the data,
// the message queues and the storage directory where the csv files will be saved.
type Server struct {
    Datastore postgresdb.Datastore
    MessageBroker rabbitmq.MessageBroker
    StorageDir string
}

func CreateServer(dbConfig postgresdb.DBConfig, MBConfig rabbitmq.MBConfig, storageDir string) (Server, error) {
    var s Server
    var err error

    // Initialize the database.
    s.Datastore, err = postgresdb.InitializeDB(dbConfig)

    if err != nil {
        return s, err
    }

    // Initialize the message broker.
    s.MessageBroker, err = rabbitmq.InitializeMB(MBConfig)

    if err != nil {
        return s, err
    }

    // Configure the storage directory.
    s.StorageDir = storageDir

    return s, nil
}
