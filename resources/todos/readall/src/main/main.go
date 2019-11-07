package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
)

// Main provides the Read All handler.
func Main(params map[string]interface{}) map[string]interface{} {
	// Setup response and headers
	res := make(map[string]interface{})
	hdrs := make(map[string]string)
	hdrs["Content-Type"] = "application/json"
	res["headers"] = hdrs

	// Get the Cloudant URL
	cloudantURL, ok := params["cloudanturl"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting cloudant url: %s", cloudantURL)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Get all the docs from the todos database.
	baseURL, err := createURL(params)
	if err != nil {
		errMsg := fmt.Sprintf("error reading paramters: %s", err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}
	todos, err := readAll(context.TODO(), cloudantURL, baseURL)
	if err != nil {
		errMsg := fmt.Sprintf("error reading docs from cloudant URL %s: %s", cloudantURL, err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Return the todo with the response.
	return jsonResponse(res, http.StatusOK, todos)
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
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Order     int    `json:"order"`
}

func readAll(ctx context.Context, cloudantURL, baseURL string) ([]todo, error) {
	var todos []todo

	// Connect to Clodoudant todos database.
	client, err := kivik.New(context.TODO(), "couch", cloudantURL)
	if err != nil {
		return todos, fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return todos, fmt.Errorf("error connecting to todos db: %s", err)
	}

	rows, err := db.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		return todos, fmt.Errorf("error getting all docs: %s", err)
	}

	// Loop through each row and create a todo from the doc, which is added to
	// the list of todos.
	for rows.Next() {
		var doc todoDoc
		if err := rows.ScanDoc(&doc); err != nil {
			return todos, fmt.Errorf("error scanning doc: %s", err)
		}
		todo := convertDocToTodo(doc)
		todo.URL = baseURL + doc.ID
		todos = append(todos, todo)
	}
	return todos, nil
}

func errResponse(res map[string]interface{}, code int, message string) map[string]interface{} {
	return jsonResponse(res, code, map[string]string{"error": message})
}

func jsonResponse(res map[string]interface{}, code int, data interface{}) map[string]interface{} {
	content, err := json.Marshal(data)
	if err != nil {
		return errResponse(res, http.StatusInternalServerError, err.Error())
	}
	if string(content) == "null" {
		content = []byte("[]")
	}
	res["statusCode"] = code
	res["body"] = content
	return res
}

func convertDocToTodo(doc todoDoc) todo {
	return todo{
		Title:     doc.Title,
		Completed: doc.Completed,
		Order:     doc.Order,
	}
}

func createURL(params map[string]interface{}) (string, error) {
	host, ok := params["ibmcloudhost"].(string)
	if !ok {
		return "", fmt.Errorf("error getting ibm cloud host: %s", host)
	}
	return host + params["__ow_path"].(string), nil
}
