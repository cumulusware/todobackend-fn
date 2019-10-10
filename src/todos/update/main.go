package main

import "net/http"

// Main provides the Update handler.
func Main(params map[string]interface{}) map[string]interface{} {
	msg := make(map[string]interface{})
	hdrs := make(map[string]string)
	msg["body"] = "Update one todo"
	msg["statusCode"] = http.StatusOK
	hdrs["Content-Type"] = "application/json"
	msg["headers"] = hdrs
	return msg
}
