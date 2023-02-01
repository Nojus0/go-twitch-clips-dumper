package main

type ClipsPayload struct {
	OperationName string            `json:"operationName"`
	Variables     Variables         `json:"variables"`
	Extensions    RequestExtensions `json:"extensions"`
}
type Criteria struct {
	Filter string `json:"filter"`
}
type Variables struct {
	Login    string   `json:"login"`
	Limit    int      `json:"limit"`
	Criteria Criteria `json:"criteria"`
	Cursor   string   `json:"cursor"`
}
type PersistedQuery struct {
	Version    int    `json:"version"`
	Sha256Hash string `json:"sha256Hash"`
}
type RequestExtensions struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}
