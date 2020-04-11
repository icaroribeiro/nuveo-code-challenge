package handlers_test

import (
    "github.com/gorilla/mux"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/router"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/server"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/joho/godotenv"
    "log"
    "os"
    "testing"
)

var s server.Server

var r *mux.Router

func init() {
    var err error

    err = godotenv.Load("../.test.env")

    if err != nil {
        log.Fatalf("Failed to load the ../.test.env file: %s", err.Error())
    }
}

// It serves as a wrapper around the testMain function that allows to defer other functions.
// At the end, it finally passes the returned exit code to os.Exit().
func TestMain(m *testing.M) {
    var exitVal int

    // Before the tests.
    utils.InitializeRandomization()

    exitVal = testMain(m)

    // After the tests.
    defer s.Datastore.Close()

    os.Exit(exitVal)
}

// It configures the settings before running the tests. It returns an integer denoting an exit code to be used 
// in the TestMain function. In the case, if the exit code is 0 it denotes success while all other codes denote failure.
func testMain(m *testing.M) int {
    var dbUsername string
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
    var err error

    // Get the database environment variables.
    dbUsername = os.Getenv("TEST_DB_USERNAME")

    if dbUsername == "" {
        log.Fatal("Failed to read the TEST_DB_USERNAME environment variable: it isn't set")
    }

    dbPassword = os.Getenv("TEST_DB_PASSWORD")

    if dbPassword == "" {
        log.Fatal("Failed to read the TEST_DB_PASSWORD environment variable: it isn't set")
    }

    dbHost = os.Getenv("TEST_DB_HOST")

    if dbHost == "" {
        log.Fatal("Failed to read the TEST_DB_HOST environment variable: it isn't set")
    }

    dbPort = os.Getenv("TEST_DB_PORT")

    if dbPort == "" {
        log.Fatal("Failed to read the TEST_DB_PORT environment variable: it isn't set")
    }

    dbName = os.Getenv("TEST_DB_NAME")

    if dbName == "" {
        log.Fatal("Failed to read the TEST_DB_NAME environment variable: it isn't set")
    }

    dbConfig = postgresdb.DBConfig{
        Username: dbUsername,
        Password: dbPassword,
        Host:     dbHost,
        Port:     dbPort,
        Name:     dbName,
    }

    // Get the message broker environment variables.
    mbUsername = os.Getenv("TEST_MB_USERNAME")

    if mbUsername == "" {
        log.Fatal("Failed to read the TEST_MB_USERNAME environment variable: it isn't set")
    }

    mbPassword = os.Getenv("TEST_MB_PASSWORD")

    if mbPassword == "" {
        log.Fatal("Failed to read the TEST_MB_PASSWORD environment variable: it isn't set")
    }

    mbHost = os.Getenv("TEST_MB_HOST")

    if mbHost == "" {
        log.Fatal("Failed to read the TEST_MB_HOST environment variable: it isn't set")
    }

    mbPort = os.Getenv("TEST_MB_PORT")

    if mbPort == "" {
        log.Fatal("Failed to read the TEST_MB_PORT environment variable: it isn't set")
    }

    mbName = os.Getenv("TEST_MB_NAME")

    if mbName == "" {
        log.Fatal("Failed to read the TEST_MB_NAME environment variable: it isn't set")
    }

    mbConfig = rabbitmq.MBConfig{
        Username: mbUsername,
        Password: mbPassword,
        Host:     mbHost,
        Port:     mbPort,
        Name:     mbName,
    }

    // Get the storage environment variable.
    storageDir = os.Getenv("STORAGE_DIR")

    if storageDir == "" {
        log.Fatal("Failed to read the STORAGE_DIR environment variable: it isn't set")
    }

    s, err = server.CreateServer(dbConfig, mbConfig, storageDir)

    if err != nil {
        log.Printf("Failed to configure the server: %s", err.Error())
        return 1
    }

    r = router.CreateRouter(&s)

    return m.Run()
}
