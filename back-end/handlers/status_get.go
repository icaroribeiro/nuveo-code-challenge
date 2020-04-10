package handlers

import (
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "net/http"
)

func GetStatus() http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var code int
        code = http.StatusOK
        utils.Respond(w, code)
    })
}
