// handlers.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Student represents the student model.
type Student struct {
	ID   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,min=1"`
}

// CreateStudentHandler creates a new student.
func CreateStudentHandler(studentsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var newStudent Student
		err := json.NewDecoder(r.Body).Decode(&newStudent)
		if err != nil {
			http.Error(w, INVALID_REQUEST_BODY, http.StatusBadRequest)
			return
		}

		// Validate request body
		if err := validate.Struct(newStudent); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Insert the student into MongoDB
		result, err := studentsCollection.InsertOne(context.Background(), newStudent)
		if err != nil {
			http.Error(w, FAILED_TO_INSERT_STUDENT, http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set(HEADER_CONTENT_TYPE, HEADER_CONTENT)
		w.WriteHeader(http.StatusCreated)

		// Send the inserted student ID in the response
		json.NewEncoder(w).Encode(result.InsertedID)
	}
}

// ReadStudentHandler retrieves a student by ID.
func ReadStudentHandler(studentsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		studentID := params["id"]

		// Query the student by ID from MongoDB
		var student Student
		err := studentsCollection.FindOne(context.Background(), bson.M{"_id": studentID}).Decode(&student)
		if err != nil {
			http.Error(w, STUDENT_NOT_FOUND, http.StatusNotFound)
			return
		}

		// Set response headers
		w.Header().Set(HEADER_CONTENT_TYPE, HEADER_CONTENT)

		// Send the student details in the response
		json.NewEncoder(w).Encode(student)
	}
}

// UpdateStudentHandler updates a student by ID.
func UpdateStudentHandler(studentsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		studentID := params["id"]

		// Parse request body
		var updatedStudent Student
		err := json.NewDecoder(r.Body).Decode(&updatedStudent)
		if err != nil {
			http.Error(w, INVALID_REQUEST_BODY, http.StatusBadRequest)
			return
		}

		// Validate request body
		if err := validate.Struct(updatedStudent); err != nil {
			http.Error(w, fmt.Sprintf("Validation error: %s", err.Error()), http.StatusBadRequest)
			return
		}

		// Update the student in MongoDB
		result, err := studentsCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": studentID},
			bson.D{{Key: "$set", Value: updatedStudent}},
		)
		if err != nil {
			http.Error(w, FAILED_TO_UPDATE_STUDENT, http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set(HEADER_CONTENT_TYPE, HEADER_CONTENT)

		// Send the number of modified documents in the response
		json.NewEncoder(w).Encode(result.ModifiedCount)
	}
}

// DeleteStudentHandler deletes a student by ID.
func DeleteStudentHandler(studentsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		studentID := params["id"]

		// Delete the student from MongoDB
		result, err := studentsCollection.DeleteOne(context.Background(), bson.M{"_id": studentID})
		if err != nil {
			http.Error(w, FAILED_TO_DELETE_STUDENT, http.StatusInternalServerError)
			return
		}

		// Set response headers
		w.Header().Set(HEADER_CONTENT_TYPE, HEADER_CONTENT)

		// Send the number of deleted documents in the response
		json.NewEncoder(w).Encode(result.DeletedCount)
	}
}
