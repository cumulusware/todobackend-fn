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

// Main provides the Delete handler.
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

	// Get the ID from the parameters
	path, ok := params["__ow_path"].(string)
	if !ok {
		errMsg := fmt.Sprintf("error getting path from parameters: %v", params)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}
	id := idFromPath(path)

	// Delete one doc from the todos database using given the ID.
	err := deleteByID(context.TODO(), cloudantURL, id)
	if err != nil {
		errMsg := fmt.Sprintf("error deleting doc id %s from cloudant URL %s: %s", id, cloudantURL, err)
		return errResponse(res, http.StatusInternalServerError, errMsg)
	}
	return jsonResponse(res, http.StatusNoContent, "")
}

type todoDoc struct {
	ID      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"`
	Deleted bool   `json:"_deleted,omitempty"`
	Title   string `json:"title"`
}

func deleteByID(ctx context.Context, url, id string) error {
	// Connect to Clodoudant todos database.
	client, err := kivik.New(context.TODO(), "couch", url)
	if err != nil {
		return fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return fmt.Errorf("error connecting to todos db: %s", err)
	}

	// Need the revision of the doc in order to delete.
	rev, err := db.Rev(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting rev of doc id %s: %s", id, err)
	}

	// Delete the doc
	_, err = db.Delete(ctx, id, rev, nil)
	if err != nil {
		return fmt.Errorf("error deleting doc ID %s: %s", id, err)
	}
	return nil
}

func errResponse(res map[string]interface{}, code int, message string) map[string]interface{} {
	return jsonResponse(res, code, map[string]string{"error": message})
}

func jsonResponse(res map[string]interface{}, code int, data interface{}) map[string]interface{} {
	var content interface{}
	var err error
	if data == "" {
		content = []byte("[]")
	} else {
		content, err = json.Marshal(data)
		if err != nil {
			return errResponse(res, http.StatusInternalServerError, err.Error())
		}
	}
	res["statusCode"] = code
	res["body"] = content
	res["body"] = "foo"
	return res
}

func idFromPath(path string) string {
	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}
