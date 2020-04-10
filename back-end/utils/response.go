package utils

import (
    "encoding/json"
    "net/http"
)

func Respond(w http.ResponseWriter, code int) {
    w.WriteHeader(code)
}

// It generates a JSON response along with the suitable header and the status code given.
func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
    var err error
    var response []byte

    response, err = json.Marshal(payload)

    if err != nil {
        w.Write([]byte("Invalid JSON response"))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}