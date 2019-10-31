package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
)

// Main provides the Create handler.
func Main(params map[string]interface{}) map[string]interface{} {
	// Setup response and headers
	res := make(map[string]interface{})
	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	res["headers"] = hdrs
	res["statusCode"] = http.StatusOK

	// Get the Cloudant URL
	cloudantURL, ok := params["cloudanturl"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting cloudant url: %s", cloudantURL)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Read the body of the request to create the todo
	title, ok := params["title"].(string)
	if !ok {
		title = "Foo"
	}
	todo := todo{
		Title: title,
	}

	// Save the todo to Cloudant
	_, err := create(context.TODO(), cloudantURL, &todo)
	if err != nil {
		errMsg := fmt.Sprintf("error saving todo to cloudant URL %s: %s", cloudantURL, err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Return the todo with the response.
	res["body"] = todo
	return res
}

type todo struct {
	Title string `json:"title"`
}

func create(ctx context.Context, url string, todo *todo) (string, error) {
	client, err := kivik.New(context.TODO(), "couch", url)
	if err != nil {
		return "", fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return "", fmt.Errorf("error connecting to todos db: %s", err)
	}

	docID, _, err := db.CreateDoc(ctx, todo)
	if err != nil {
		return "", fmt.Errorf("error creating doc: %s", err)
	}
	return docID, nil
}

func errResponse(res map[string]interface{}, code int, message string) map[string]interface{} {
	return jsonResponse(res, code, map[string]string{"error": message})
}

func jsonResponse(res map[string]interface{}, code int, data interface{}) map[string]interface{} {
	content, err := json.Marshal(data)
	if string(content) == "null" {
		content = []byte("[]")
	}
	if err != nil {
		errResponse(res, http.StatusInternalServerError, err.Error())
	}
	res["statusCode"] = code
	res["body"] = content
	return res
}
