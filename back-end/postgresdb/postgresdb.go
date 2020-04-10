package postgresdb

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

// The DBConfig stores all the parameters to configure the database settings.
type DBConfig struct {
    Username string
    Password string
    Host     string
    Port     string
    Name     string
}

// The Datastore groups all the variables necessary to connect and work with data 
// by means of a collection of statements used to interface with our backing database.
type Datastore struct {
    DB *sql.DB
    Stmts map[string]*sql.Stmt
}

func InitializeDB(dbConfig DBConfig) (Datastore, error) {
    var connString string
    var db *sql.DB
    var err error
    var unpreparedStmts map[string]string
    var preparedStmts map[string]*sql.Stmt

    // Set up the connection string of the database.
    connString = SetUpConnString(dbConfig)

    db, err = sql.Open("postgres", connString)

    if err != nil {
        return Datastore{}, fmt.Errorf("it wasn't possible to open the database: %s", err.Error())
    }

    err = db.Ping()

    if err != nil {
        return Datastore{}, fmt.Errorf("it wasn't possible to connect to the database: %s", err.Error())
    }

    unpreparedStmts = make(map[string]string)

    ConfigureStatements(unpreparedStmts)

    preparedStmts = make(map[string]*sql.Stmt)

    err = PrepareStatements(db, preparedStmts, unpreparedStmts)

    if err != nil {
        return Datastore{}, fmt.Errorf("it wasn't possible to prepare the database statements: %s", err.Error())
    }

    return Datastore{DB: db, Stmts: preparedStmts}, nil
}

// It builds the connection string of the database.
func SetUpConnString(dbConfig DBConfig) string {
    var connString string

    connString = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
                        dbConfig.Username,
                        dbConfig.Password,
                        dbConfig.Host,
                        dbConfig.Port,
                        dbConfig.Name)

    return connString
}

// It configures all statements of CRUD operations.
func ConfigureStatements(unpreparedStmts map[string]string) {
    // Workflows.
    AddCreateWorkflowStatement(unpreparedStmts)
    AddGetAllWorkflowsStatement(unpreparedStmts)
    AddGetWorkflowStatement(unpreparedStmts)
    AddUpdateWorkflowStatement(unpreparedStmts)
}

// It prepares each query statement on the database to verify syntaxes.
// If one of them fails, it returns an error.
func PrepareStatements(db *sql.DB, preparedStmts map[string]*sql.Stmt, unpreparedStmts map[string]string) error {
    var key string
    var value string
    var err error
    var stmt *sql.Stmt

    for key, value = range unpreparedStmts {
        stmt, err = db.Prepare(value)

        if err != nil {
            return err
        }

        preparedStmts[key] = stmt
    }

    return nil
}

// It closes the database connection, releasing any open resources.
func (d *Datastore) Close() error {
    return d.DB.Close()
}
