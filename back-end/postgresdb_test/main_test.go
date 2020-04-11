package postgresdb_test

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "github.com/joho/godotenv"
    "log"
    "os"
    "testing"
)

var datastore postgresdb.Datastore

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
    defer datastore.Close()

    os.Exit(exitVal)
}

func testMain(m *testing.M) int {
    var dbUsername string
    var dbPassword string
    var dbHost string
    var dbPort string
    var dbName string
    var dbConfig postgresdb.DBConfig
    var err error

    dbUsername = os.Getenv("TEST_DB_USERNAME")

    if dbUsername == "" {
        log.Print("Failed to read the TEST_DB_USERNAME environment variable: it isn't set")
        return 1
    }

    dbPassword = os.Getenv("TEST_DB_PASSWORD")

    if dbPassword == "" {
        log.Print("Failed to read the TEST_DB_PASSWORD environment variable: it isn't set")
        return 1
    }

    dbHost = os.Getenv("TEST_DB_HOST")

    if dbHost == "" {
        log.Print("Failed to read the TEST_DB_HOST environment variable: it isn't set")
        return 1
    }

    dbPort = os.Getenv("TEST_DB_PORT")

    if dbPort == "" {
        log.Print("Failed to read the TEST_DB_PORT environment variable: it isn't set")
        return 1
    }

    dbName = os.Getenv("TEST_DB_NAME")

    if dbName == "" {
        log.Print("Failed to read the TEST_DB_NAME environment variable: it isn't set")
        return 1
    }

    log.Printf(dbName)

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

    return m.Run()
}
