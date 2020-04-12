package rabbitmq_test

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/postgresdb"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/rabbitmq"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "log"
    "os"
    "testing"
)

var envVariablesMap map[string]string

var datastore postgresdb.Datastore

var messageBroker rabbitmq.MessageBroker

func init() {
    var filenames []string
    var err error

    filenames = []string{"../.test.env"}

    envVariablesMap = make(map[string]string)

    envVariablesMap["TEST_DB_USERNAME"] = ""
    envVariablesMap["TEST_DB_PASSWORD"] = ""
    envVariablesMap["TEST_DB_HOST"] = ""
    envVariablesMap["TEST_DB_PORT"] = ""
    envVariablesMap["TEST_DB_NAME"] = ""

    envVariablesMap["TEST_MB_USERNAME"] = ""
    envVariablesMap["TEST_MB_PASSWORD"] = ""
    envVariablesMap["TEST_MB_HOST"] = ""
    envVariablesMap["TEST_MB_PORT"] = ""
    envVariablesMap["TEST_MB_NAME"] = ""

    err = utils.GetEnvVariables(filenames, envVariablesMap)

    if err != nil {
        log.Fatal(err.Error())
    }
}

func TestMain(m *testing.M) {
    var exitVal int

    // Before the tests.
    utils.InitializeRandomization()

    exitVal = testMain(m)

    // After the tests.
    defer messageBroker.Close()

    os.Exit(exitVal)
}

func testMain(m *testing.M) int {
    var dbConfig postgresdb.DBConfig
    var err error
    var mbConfig rabbitmq.MBConfig

    dbConfig = postgresdb.DBConfig{
        Username: envVariablesMap["TEST_DB_USERNAME"],
        Password: envVariablesMap["TEST_DB_PASSWORD"],
        Host:     envVariablesMap["TEST_DB_HOST"],
        Port:     envVariablesMap["TEST_DB_PORT"],
        Name:     envVariablesMap["TEST_DB_NAME"],
    }

    datastore, err = postgresdb.InitializeDB(dbConfig)

    if err != nil {
        log.Printf("Failed to configure the database: %s", err.Error())
        return 1
    }

    mbConfig = rabbitmq.MBConfig{
        Username: envVariablesMap["TEST_MB_USERNAME"],
        Password: envVariablesMap["TEST_MB_PASSWORD"],
        Host:     envVariablesMap["TEST_MB_HOST"],
        Port:     envVariablesMap["TEST_MB_PORT"],
        Name:     envVariablesMap["TEST_MB_NAME"],
    }

    messageBroker, err = rabbitmq.InitializeMB(mbConfig)

    if err != nil {
        log.Printf("Failed to configure the message broker: %s", err.Error())
        return 1
    }

    return m.Run()
}
