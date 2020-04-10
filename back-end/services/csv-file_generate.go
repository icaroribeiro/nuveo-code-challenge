package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/yukithm/json2csv"
    "os"
)

func GenerateCSVFile(jsonBytes []byte, storageDir string, workflowId string) error {
    var reader *bytes.Reader
    var decoder *json.Decoder
    var err error
    var jsonObj interface{}
    var results []json2csv.KeyValue
    var file *os.File
    var csv *json2csv.CSVWriter

    reader = bytes.NewReader(jsonBytes)
    
    decoder = json.NewDecoder(reader)
    
    // It causes the decoder to unmarshal a number 
    // into an interface{} as a Number instead of as a float64.
    decoder.UseNumber()

    err = decoder.Decode(&jsonObj)

    if err != nil {
        return err
    }

    // It represents an array of key(path)/value map, same as []map[string]interface{}.
    results, err = json2csv.JSON2CSV(jsonObj)

    if err != nil {
        return err
    }

    _, err = os.Stat(storageDir)
 
	if os.IsNotExist(err) {
        err = os.MkdirAll(storageDir, 0755)
        
		if err != nil {
			return err
		}
    }

    file, err = os.Create(fmt.Sprintf("%s/%s.csv", storageDir, workflowId))

    if err != nil {
        return err
    }

    defer file.Close()

    csv = json2csv.NewCSVWriter(file)

    // Set up the dot-bracket header style which uses square brackets for array indexes.
    csv.HeaderStyle = json2csv.DotBracketStyle

    err = csv.WriteCSV(results)

    if err != nil {
        return err
    }

    return nil
}
