package main

import "net/http"

// Main provides the Delete All handler.
func Main(params map[string]interface{}) map[string]interface{} {
	msg := make(map[string]interface{})
	hdrs := make(map[string]string)
	msg["body"] = []byte("[]")
	msg["statusCode"] = http.StatusNoContent
	hdrs["Content-Type"] = "application/json"
	msg["headers"] = hdrs
	return msg
}
