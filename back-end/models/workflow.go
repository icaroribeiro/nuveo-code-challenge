package models

type Workflow struct {
    ID     string   `json:"id,omitempty"`
    Status string   `json:"status"`
    Data   DataMap  `json:"data"`
    Steps  []string `json:"steps"`
}
