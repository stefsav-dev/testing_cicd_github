package main

import (
	"latihan_devops/internal/database"
	"latihan_devops/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Database connection failed : ", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	userHandler := handlers.NewUserHandler(db)

	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUSer).Methods("POST")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	log.Println("Server running on port 6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}
