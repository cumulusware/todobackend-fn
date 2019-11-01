package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb"
)

// Main provides the Delete All handler.
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

	// Delete all the docs from the todos database.
	err := deleteAll(context.TODO(), cloudantURL)
	if err != nil {
		errMsg := fmt.Sprintf("error reading docs from cloudant URL %s: %s", cloudantURL, err)
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

func deleteAll(ctx context.Context, url string) error {
	// Connect to Clodoudant todos database.
	client, err := kivik.New(context.TODO(), "couch", url)
	if err != nil {
		return fmt.Errorf("error opening couchdb: %s", err)
	}
	db, err := client.DB(ctx, "todos")
	if err != nil {
		return fmt.Errorf("error connecting to todos db: %s", err)
	}

	// Get all docs
	var docs []todoDoc
	rows, err := db.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		return fmt.Errorf("error getting all docs: %s", err)
	}

	// Iterate through each doc and set to deleted.
	for rows.Next() {
		var doc todoDoc
		if err := rows.ScanDoc(&doc); err != nil {
			return fmt.Errorf("error scanning doc: %s", err)
		}
		doc.Deleted = true
		docs = append(docs, doc)
	}

	// Bulk update all docs to be deleted.
	time.Sleep(300 * time.Millisecond) // Added for IBM Cloud rate limit on lite plan.
	_, err = db.BulkDocs(ctx, docs, nil)
	return err
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
