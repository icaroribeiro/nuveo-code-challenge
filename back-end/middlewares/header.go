package middlewares

import (
    "fmt"
    "net/http"
    "github.com/icaroribeiro/nuveo-code-challenge/back-end/utils"
    "strings"
)

// It validates the request header fields with the ones expected from server side based on a key-value pair scheme.
func ValidateRequestHeaderFields(fields map[string]string) Adapter {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            var key string
            var value string
            var code int

            for key, value = range fields {            
                if !strings.Contains(r.Header.Get(key), value) {
                    code = http.StatusBadRequest 
                    utils.RespondWithJson(w, code, map[string]string{"error": 
                        fmt.Sprintf("the request header key %s must be related with the value %s", key, value)})
                    return
                }
            }

            // Call the next handler which can be another middleware in the chain or the final handler.
            next.ServeHTTP(w, r)
        }
    }
}