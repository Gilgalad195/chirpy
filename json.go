package main

type parameters struct {
	Body string `json:"body"`
}

type cleaned struct {
	CleanedBody string `json:"cleaned_body"`
}
type errorResp struct {
	Error string `json:"error"`
}

// type validResp struct {
// 	Valid bool `json:"valid"`
// }
