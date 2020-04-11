package rabbitmq_test

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/joho/godotenv"
    "log"
    "os"
    "testing"
)

var datastore postgresdb.Datastore

var messageQueue rabbitmq.MessageBroker

func init() {
    var err error

    err = godotenv.Load("../.test.env")

    if err != nil {
        log.Fatalf("Failed to load the ../.test.env file: %s", err.Error())
    }
}

func TestMain(m *testing.M) {
    var exitVal int

    // Before the tests.
    utils.InitializeRandomization()

    exitVal = testMain(m)

    // After the tests.
    defer messageQueue.Close()

    os.Exit(exitVal)
}

func testMain(m *testing.M) int {
    var dbUsername string
    var isSet bool
    var dbPassword string
    var dbHost string
    var dbPort string
    var dbName string
    var dbConfig postgresdb.DBConfig
    var err error
    var mbUsername string
    var mbPassword string
    var mbHost string
    var mbPort string
    var mbName string
    var mbConfig rabbitmq.MBConfig

    dbUsername, isSet = os.LookupEnv("TEST_DB_USERNAME")

    if !isSet {
        log.Fatal("Failed to read the TEST_DB_USERNAME environment variable: it isn't set")
    }

    dbPassword, isSet = os.LookupEnv("TEST_DB_PASSWORD")

    if !isSet {
        log.Fatal("Failed to read the TEST_DB_PASSWORD environment variable: it isn't set")
    }

    dbHost, isSet = os.LookupEnv("TEST_DB_HOST")

    if !isSet {
        log.Fatal("Failed to read the TEST_DB_HOST environment variable: it isn't set")
    }

    dbPort, isSet = os.LookupEnv("TEST_DB_PORT")

    if !isSet {
        log.Fatal("Failed to read the TEST_DB_PORT environment variable: it isn't set")
    }

    dbName, isSet = os.LookupEnv("TEST_DB_NAME")

    if !isSet {
        log.Fatal("Failed to read the TEST_DB_NAME environment variable: it isn't set")
    }

    dbConfig = postgresdb.DBConfig{
        Username: dbUsername,
        Password: dbPassword,
        Host:     dbHost,
        Port:     dbPort,
        Name:     dbName,
    }

    datastore, err = postgresdb.InitializeDB(dbConfig)

    if err != nil {
        log.Printf("Failed to configure the database: %s", err.Error())
        return 1
    }

    mbUsername, isSet = os.LookupEnv("TEST_MB_USERNAME")

    if !isSet {
        log.Fatal("Failed to read the TEST_MB_USERNAME environment variable: it isn't set")
    }

    mbPassword, isSet = os.LookupEnv("TEST_MB_PASSWORD")

    if !isSet {
        log.Fatal("Failed to read the TEST_MB_PASSWORD environment variable: it isn't set")
    }

    mbHost, isSet = os.LookupEnv("TEST_MB_HOST")

    if !isSet {
        log.Fatal("Failed to read the TEST_MB_HOST environment variable: it isn't set")
    }

    mbPort, isSet = os.LookupEnv("TEST_MB_PORT")

    if !isSet {
        log.Fatal("Failed to read the TEST_MB_PORT environment variable: it isn't set")
    }

    mbName, isSet = os.LookupEnv("TEST_MB_NAME")

    if !isSet {
        log.Fatal("Failed to read the TEST_MB_NAME environment variable: it isn't set")
    }

    mbConfig = rabbitmq.MBConfig{
        Username: mbUsername,
        Password: mbPassword,
        Host:     mbHost,
        Port:     mbPort,
        Name:     mbName,
    }

    messageQueue, err = rabbitmq.InitializeMB(mbConfig)

    if err != nil {
        log.Printf("Failed to configure the message queue: %s", err.Error())
        return 1
    }

    return m.Run()
}
