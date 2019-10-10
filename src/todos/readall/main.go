package main

import "net/http"

// Main provides the Read All handler.
func Main(params map[string]interface{}) map[string]interface{} {
	msg := make(map[string]interface{})
	hdrs := make(map[string]string)
	msg["body"] = "Read all todos"
	msg["statusCode"] = http.StatusOK
	hdrs["Content-Type"] = "application/json"
	msg["headers"] = hdrs
	return msg
}
