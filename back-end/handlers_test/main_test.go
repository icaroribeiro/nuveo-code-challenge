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
    var err error

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

    storageDir, isSet = os.LookupEnv("STORAGE_DIR")

    if !isSet {
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