package main

import (
    "fmt"
    "github.com/gorilla/mux"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/router"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/joho/godotenv"
    "log"
    "net/http"
    "os"
    "os/signal"
)

func init() {
    var err error

    // Load the variables from .env file into the system.
    err = godotenv.Load()

    if err != nil {
        log.Fatal("Failed to load the .env file: ", err.Error())
    }
}

func main() {
    var dbUsername string
    var isSet bool
    var dbPassword string
    var dbHost string
    var dbPort string
    var dbName string
    var dbConfig postgresdb.DBConfig
    var mbUsername string
    var mbPassword string
    var mbHost string
    var mbPort string
    var mbName string
    var mbConfig rabbitmq.MBConfig
    var storageDir string
    var s server.Server
    var err error
    var r *mux.Router
    var httpPort string
    var httpHost string
    var httpAddress string

    // Get the database environment variables.
    dbUsername, isSet = os.LookupEnv("DB_USERNAME")

    if !isSet {
        log.Fatal("Failed to read the DB_USERNAME environment variable: it isn't set")
    }

    dbPassword, isSet = os.LookupEnv("DB_PASSWORD")

    if !isSet {
        log.Fatal("Failed to read the DB_PASSWORD environment variable: it isn't set")
    }

    dbHost, isSet = os.LookupEnv("DB_HOST")

    if !isSet {
        log.Fatal("Failed to read the DB_HOST environment variable: it isn't set")
    }

    dbPort, isSet = os.LookupEnv("DB_PORT")

    if !isSet {
        log.Fatal("Failed to read the DB_PORT environment variable: it isn't set")
    }

    dbName, isSet = os.LookupEnv("DB_NAME")

    if !isSet {
        log.Fatal("Failed to read the DB_NAME environment variable: it isn't set")
    }

    dbConfig = postgresdb.DBConfig{
        Username: dbUsername,
        Password: dbPassword,
        Host:     dbHost,
        Port:     dbPort,
        Name:     dbName,
    }

    // Get the message broker environment variables.
    mbUsername, isSet = os.LookupEnv("MB_USERNAME")

    if !isSet {
        log.Fatal("Failed to read the MB_USERNAME environment variable: it isn't set")
    }

    mbPassword, isSet = os.LookupEnv("MB_PASSWORD")

    if !isSet {
        log.Fatal("Failed to read the MB_PASSWORD environment variable: it isn't set")
    }

    mbHost, isSet = os.LookupEnv("MB_HOST")

    if !isSet {
        log.Fatal("Failed to read the MB_HOST environment variable: it isn't set")
    }

    mbPort, isSet = os.LookupEnv("MB_PORT")

    if !isSet {
        log.Fatal("Failed to read the MB_PORT environment variable: it isn't set")
    }

    mbName, isSet = os.LookupEnv("MB_NAME")

    if !isSet {
        log.Fatal("Failed to read the MB_NAME environment variable: it isn't set")
    }

    mbConfig = rabbitmq.MBConfig{
        Username: mbUsername,
        Password: mbPassword,
        Host:     mbHost,
        Port:     mbPort,
        Name:     mbName,
    }

    // Get the storage environment variable.
    storageDir, isSet = os.LookupEnv("STORAGE_DIR")

    if !isSet {
        log.Fatal("Failed to read the STORAGE_DIR environment variable: it isn't set")
    }

    // Create the server.
    s, err = server.CreateServer(dbConfig, mbConfig, storageDir)

    if err != nil {
        log.Fatal("Failed to configure the server: ", err.Error())
    }

    // Create the router by arranging the routes.
    r = router.CreateRouter(&s)

    // Get the http server environment variables.
    httpHost, isSet = os.LookupEnv("HTTP_SERVER_HOST")

    if !isSet {
        log.Fatal("Failed to read the HTTP_SERVER_HOST environment variable: it isn't set")
    }

    httpPort, isSet = os.LookupEnv("HTTP_SERVER_PORT")

    if !isSet {
        log.Fatal("Failed to read the HTTP_SERVER_PORT environment variable: it isn't set")
    }

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
        log.Fatalf("Failed to close the message queue: %s", err.Error())
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
