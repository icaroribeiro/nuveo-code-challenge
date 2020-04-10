package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGetStatus(t *testing.T) {
    var method string
    var path string
    var request *http.Request
    var err error
    var response *httptest.ResponseRecorder
    var expectedCode int

    method = "GET"

    path = "/status"

    request, err = http.NewRequest(method, path, nil)

    if err != nil {
        t.Fatalf("Failed to create the request: %s", err.Error())
    }

    t.Logf("Request: method=%s and path=%s", method, path)

    response = httptest.NewRecorder()

    r.ServeHTTP(response, request)

    expectedCode = http.StatusOK

    if expectedCode != response.Code {
        t.Errorf("Test failed, response: code=%d and body=%+v", response.Code, response.Body)
        return
    }

    t.Logf("Test successful, response: code=%d", response.Code)
}
