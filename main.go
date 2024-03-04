// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate *validator.Validate

func main() {
	// Initialize MongoDB client and collection
	mongoClient, err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	studentsCollection := mongoClient.Database("school").Collection("students")

	// Create router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/students", CreateStudentHandler(studentsCollection)).Methods("POST")
	router.HandleFunc("/students/{id}", ReadStudentHandler(studentsCollection)).Methods("GET")
	router.HandleFunc("/students/{id}", UpdateStudentHandler(studentsCollection)).Methods("PUT")
	router.HandleFunc("/students/{id}", DeleteStudentHandler(studentsCollection)).Methods("DELETE")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initMongoDB() (*mongo.Client, error) {
	// Initialize MongoDB client (replace the connection string with your MongoDB instance)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	// Connect to MongoDB
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}
