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

// Main provides the Update handler.
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

	// Read the body of the request to create the todo
	title, ok := params["title"].(string)
	if !ok {
		title = "Foo"
	}
	todo := todo{
		Title:     title,
		Completed: false,
	}
	if completed, ok := params["completed"].(bool); ok {
		todo.Completed = completed
	}
	if orderFloat, ok := params["order"].(float64); ok {
		todo.Order = int(orderFloat)
	}

	// Get the ID from the parameters
	path, ok := params["__ow_path"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting path from parameters: %v", params)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}
	id := idFromPath(path)

	// Update the todo
	err := updateByID(context.TODO(), cloudantURL, id, todo)
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

func updateByID(ctx context.Context, cloudantURL, id string, todo todo) error {
	// Connect to Clodoudant todos database.
	client, err := kivik.New(ctx, "couch", cloudantURL)
	if err != nil {
		return fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return fmt.Errorf("error connecting to todos db: %s", err)
	}

	// Need the revision of the doc in order to update.
	rev, err := db.Rev(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting rev of doc id %s: %s", id, err)
	}

	doc := convertTodoToDoc(todo)
	doc.ID = id
	doc.Rev = rev
	_, err = db.Put(ctx, id, doc)
	if err != nil {
		return fmt.Errorf("error putting doc ID %s: %s", id, err)
	}
	return nil
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
	Order     int    `json:"order"`
}

type todoDoc struct {
	ID        string `json:"_id,omitempty"`
	Rev       string `json:"_rev,omitempty"`
	Deleted   bool   `json:"_deleted,omitempty"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

func convertTodoToDoc(todo todo) todoDoc {
	return todoDoc{
		Title:     todo.Title,
		Completed: todo.Completed,
		Order:     todo.Order,
	}
}

func createURL(params map[string]interface{}) (string, error) {
	host, ok := params["ibmcloudhost"].(string)
	if !ok {
		return "", fmt.Errorf("error getting ibm cloud host: %s", host)
	}
	return host + params["__ow_path"].(string), nil
}
