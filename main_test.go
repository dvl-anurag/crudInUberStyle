// main_test.go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCRUD(t *testing.T) {
	// Initialize MongoDB client and collection for testing
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	testDB := mongoClient.Database("testdb")
	testCollection := testDB.Collection("teststudents")

	// Create router
	router := mux.NewRouter()

	// Define routes with the testing collection
	router.HandleFunc("/students", CreateStudentHandler(testCollection)).Methods("POST")
	router.HandleFunc("/students/{id}", ReadStudentHandler(testCollection)).Methods("GET")
	router.HandleFunc("/students/{id}", UpdateStudentHandler(testCollection)).Methods("PUT")
	router.HandleFunc("/students/{id}", DeleteStudentHandler(testCollection)).Methods("DELETE")

	// Run tests
	t.Run("CreateStudent", func(t *testing.T) {
		createStudentTest(t, router)
	})
	t.Run("ReadStudent", func(t *testing.T) {
		readStudentTest(t, router)
	})
	t.Run("UpdateStudent", func(t *testing.T) {
		updateStudentTest(t, router)
	})
	t.Run("DeleteStudent", func(t *testing.T) {
		deleteStudentTest(t, router)
	})
}

func createStudentTest(t *testing.T, router *mux.Router) {
	payload := `{"name": "John Doe", "age": 25}`
	req, err := http.NewRequest("POST", "/students", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is 201 Created
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body for the inserted ID
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := result["$oid"]; !ok {
		t.Errorf("Expected inserted ID in response body, got: %v", result)
	}
}

func readStudentTest(t *testing.T, router *mux.Router) {
	// Insert a test student for reading
	payload := `{"name": "Jane Doe", "age": 30}`
	req, err := http.NewRequest("POST", "/students", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Extract the inserted ID from the response
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}
	insertedID := result["$oid"].(string)

	// Read the student using the inserted ID
	req, err = http.NewRequest("GET", "/students/"+insertedID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body for the retrieved student
	var retrievedStudent map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &retrievedStudent)
	if err != nil {
		t.Fatal(err)
	}

	if retrievedStudent["_id"] != insertedID {
		t.Errorf("Expected retrieved student ID to match inserted ID, got: %v", retrievedStudent)
	}
}

func updateStudentTest(t *testing.T, router *mux.Router) {
	// Insert a test student for updating
	payload := `{"name": "Jack Doe", "age": 35}`
	req, err := http.NewRequest("POST", "/students", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Extract the inserted ID from the response
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}
	insertedID := result["$oid"].(string)

	// Update the student using the inserted ID
	updatePayload := `{"name": "Updated Name", "age": 40}`
	req, err = http.NewRequest("PUT", "/students/"+insertedID, strings.NewReader(updatePayload))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body for the number of modified documents
	var modifiedCount int64
	err = json.Unmarshal(rr.Body.Bytes(), &modifiedCount)
	if err != nil {
		t.Fatal(err)
	}

	if modifiedCount != 1 {
		t.Errorf("Expected one modified document, got: %v", modifiedCount)
	}
}

func deleteStudentTest(t *testing.T, router *mux.Router) {
	// Insert a test student for deleting
	payload := `{"name": "Deleted Student", "age": 50}`
	req, err := http.NewRequest("POST", "/students", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Extract the inserted ID from the response
	var result map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}
	insertedID := result["$oid"].(string)

	// Delete the student using the inserted ID
	req, err = http.NewRequest("DELETE", "/students/"+insertedID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body for the number of deleted documents
	var deletedCount int64
	err = json.Unmarshal(rr.Body.Bytes(), &deletedCount)
	if err != nil {
		t.Fatal(err)
	}

	if deletedCount != 1 {
		t.Errorf("Expected one deleted document, got: %v", deletedCount)
	}
}
