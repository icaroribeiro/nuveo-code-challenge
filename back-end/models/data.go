package models

import (
    "encoding/json"
    "database/sql/driver"
    "fmt"
)

/*
When querying the JSONB column of the Postgres database, the lib/pq package driver will return a bytestring.
Thus, the most convenient way to work with JSONB coming from a database would be in the form of a map[string]interface{}.
Luckely, the database/sql package has 2 built-in interfaces that we can use to create our own database compatible type: 
    i) sql.Scanner
    ii) driver.Value
Basically, if we have a type that implements these 2 interfaces, we can use that type with the database/sql package.
*/

// First, it is necessary to create the type for the data field. 
// Afterwards, it is required to implement the interface.
type DataMap map[string]interface{}

func (d DataMap) IsEmpty() bool {
    if (d == nil) {
        return true
    }

    return false
}

// To satisfy the driver.Value interface it must be implemented the Value method
// which must transform our type to a database driver compatible type. 
// In the case, it will marshall the map to JSONB data (= []byte).
func (d DataMap) Value() (driver.Value, error) {
    var bytes []byte
    var err error

    bytes, err = json.Marshal(d)
    
    return bytes, err
}

// To comply with the sql.Scanner interface it must be implemented the Scan method
// which must take the raw data that comes from the database and transform it to our new type.
// In the case, the database will return JSONB data (= []byte) that must be transformed to our type.
// (In other words, the procedure is the reverse of what is done with driver.Value)
func (d *DataMap) Scan(src interface{}) error {
    var source []byte
    var isOK bool
    var i interface{}
    var err error
    
    source, isOK = src.([]byte)

    if !isOK {
        return fmt.Errorf("type assertion .([]byte) failed")
    }

    err = json.Unmarshal(source, &i)

    if err != nil {
        return err
    }

    *d, isOK = i.(map[string]interface{})

    if !isOK {
        return fmt.Errorf("type assertion .(map[string]interface{}) failed")
    }

    return nil
}
