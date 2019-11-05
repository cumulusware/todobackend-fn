package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
)

// Main provides the Read handler.
func Main(params map[string]interface{}) map[string]interface{} {
	// Setup response and headers
	res := make(map[string]interface{})
	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	res["headers"] = hdrs

	// Get the Cloudant URL
	cloudantURL, ok := params["cloudanturl"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting cloudant url from parameters: %s", cloudantURL)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Get one todo by ID.
	path, ok := params["__ow_path"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting path from parameters: %v", params)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}
	id := idFromPath(path)

	todo, err := getByID(context.TODO(), id, cloudantURL)
	if err != nil {
		errMsg := fmt.Sprintf("error reading doc id %s from cloudant URL %s: %s", id, cloudantURL, err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Create the URL for the todo.
	todo.URL, err = createURL(params)
	if err != nil {
		return errResponse(res, http.StatusInternalServerError, err.Error())
	}

	// Return the todo with the response.
	return jsonResponse(res, http.StatusOK, todo)
}

func getByID(ctx context.Context, id string, url string) (todo, error) {
	var todo todo

	// Connect to Clodoudant todos database.
	client, err := kivik.New(context.TODO(), "couch", url)
	if err != nil {
		return todo, fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return todo, fmt.Errorf("error connecting to todos db: %s", err)
	}

	// Read by ID one doc.
	row, err := db.Get(ctx, id, nil)
	if err != nil {
		return todo, fmt.Errorf("error getting doc with ID %s: %s", id, err)
	}
	var doc todoDoc
	if err := row.ScanDoc(&doc); err != nil {
		return todo, fmt.Errorf("error scanning doc: %s", err)
	}
	todo = convertDocToTodo(doc)
	return todo, nil
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

func idFromPath(path string) string {
	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}

type todo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	URL       string `json:"url"`
}

type todoDoc struct {
	ID        string `json:"_id,omitempty"`
	Rev       string `json:"_rev,omitempty"`
	Deleted   bool   `json:"_deleted,omitempty"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func convertDocToTodo(doc todoDoc) todo {
	return todo{
		Title:     doc.Title,
		Completed: doc.Completed,
	}
}

func createURL(params map[string]interface{}) (string, error) {
	host, ok := params["ibmcloudhost"].(string)
	if !ok {
		return "", fmt.Errorf("error getting ibm cloud host: %s", host)
	}
	return host + params["__ow_path"].(string), nil
}
