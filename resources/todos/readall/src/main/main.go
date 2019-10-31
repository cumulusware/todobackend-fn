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
	res["statusCode"] = http.StatusOK

	// Get the Cloudant URL
	cloudantURL, ok := params["cloudanturl"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting cloudant url: %s", cloudantURL)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Get all the docs from the todos database.
	todos, err := readAll(context.TODO(), cloudantURL)
	if err != nil {
		errMsg := fmt.Sprintf("error reading docs from cloudant URL %s: %s", cloudantURL, err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}

	// Return the todo with the response.
	res["body"] = todos
	return res
}

type todo struct {
	Title string `json:"title"`
}

type todoDoc struct {
	ID    string `json:"_id,omitempty"`
	Rev   string `json:"_rev,omitempty"`
	Title string `json:"title"`
}

func readAll(ctx context.Context, url string) ([]todo, error) {
	var todos []todo
	client, err := kivik.New(context.TODO(), "couch", url)
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
		todos = append(todos, todo)
	}
	return todos, nil
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

func convertDocToTodo(doc todoDoc) todo {
	return todo{
		Title: doc.Title,
	}
}
