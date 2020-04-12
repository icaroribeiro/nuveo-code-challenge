package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/router"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "log"
    "net/http"
    "os"
    "os/signal"
)

var envVariablesMap map[string]string

func init() {
    var filenames []string
    var err error

    envVariablesMap = make(map[string]string)

    // The environment variables related to the database settings.
    envVariablesMap["DB_USERNAME"] = ""
    envVariablesMap["DB_PASSWORD"] = ""
    envVariablesMap["DB_HOST"] = ""
    envVariablesMap["DB_PORT"] = ""
    envVariablesMap["DB_NAME"] = ""

    // The environment variables of the message broker settings.
    envVariablesMap["MB_USERNAME"] = ""
    envVariablesMap["MB_PASSWORD"] = ""
    envVariablesMap["MB_HOST"] = ""
    envVariablesMap["MB_PORT"] = ""
    envVariablesMap["MB_NAME"] = ""

    // The environment variable of the storage settings.
    envVariablesMap["STORAGE_DIR"] = ""

    // The environment variables related to the HTTP server.
    envVariablesMap["HTTP_SERVER_HOST"] = ""
    envVariablesMap["HTTP_SERVER_PORT"] = ""

    // The environment file from where the variables will be loaded.
    filenames = []string{"./.env"}

    err = utils.GetEnvVariables(filenames, envVariablesMap)

    if err != nil {
        log.Fatal(err.Error())
    }
}

func main() {
    var dbConfig postgresdb.DBConfig
    var mbConfig rabbitmq.MBConfig
    var storageDir string
    var s server.Server
    var err error
    var r *mux.Router
    var httpPort string
    var httpHost string
    var httpAddress string

    dbConfig = postgresdb.DBConfig{
        Username: envVariablesMap["DB_USERNAME"],
        Password: envVariablesMap["DB_PASSWORD"],
        Host:     envVariablesMap["DB_HOST"],
        Port:     envVariablesMap["DB_PORT"],
        Name:     envVariablesMap["DB_NAME"],
    }

    mbConfig = rabbitmq.MBConfig{
        Username: envVariablesMap["MB_USERNAME"],
        Password: envVariablesMap["MB_PASSWORD"],
        Host:     envVariablesMap["MB_HOST"],
        Port:     envVariablesMap["MB_PORT"],
        Name:     envVariablesMap["MB_NAME"],
    }

    storageDir = envVariablesMap["STORAGE_DIR"]

    // Create the server.
    s, err = server.CreateServer(dbConfig, mbConfig, storageDir)

    if err != nil {
        log.Fatal("Failed to configure the server: ", err.Error())
    }

    // Create the router by arranging the routes.
    r = router.CreateRouter(&s)

    httpHost = envVariablesMap["HTTP_SERVER_HOST"]
    httpPort = envVariablesMap["HTTP_SERVER_PORT"]

    httpAddress = fmt.Sprintf("%s:%s", httpHost, httpPort)

    log.Printf("Starting the HTTP server connection on %s", httpAddress)

    go func() {
        err = http.ListenAndServe(httpAddress, r)

        if err != nil {
            log.Fatalf("Failed to start the HTTP server connection to %s: %s", httpAddress, err.Error())
        }
    }()

    // Graceful disconnect.
    WaitForShutdown()

    err = s.Datastore.Close()

    if err != nil {
        log.Fatalf("Failed to close the database: %s", err.Error())
    }

    s.MessageBroker.Close()

    if err != nil {
        log.Fatalf("Failed to close the message broker: %s", err.Error())
    }

    log.Println("Done")
}

func WaitForShutdown() {
    var interruptChan chan os.Signal

    // Create a channel to receive OS signals.
    interruptChan = make(chan os.Signal)

    // Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
    // ignoring other incoming signals.
    signal.Notify(interruptChan, os.Interrupt)

    // Block the main routine so that to keep it running until a signal is received.
    // If the main routine is shut down, the child one that is serving the server will shut down as well.
    <-interruptChan

    log.Println("Shutting down the server...")
}
