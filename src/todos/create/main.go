package main

import "net/http"

// Main provides the Create handler.
func Main(params map[string]interface{}) map[string]interface{} {
	msg := make(map[string]interface{})
	hdrs := make(map[string]string)
	todo := todo{
		Title: "a todo",
	}
	msg["body"] = todo
	msg["statusCode"] = http.StatusOK
	hdrs["Content-Type"] = "application/json"
	msg["headers"] = hdrs
	return msg
}

type todo struct {
	Title string `json:"title"`
}
