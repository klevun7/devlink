package main

import (
	"devlink/api/models"      
	"devlink/api/notifications" 
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	log.Println("Migrating database schema...")
	db.AutoMigrate(&models.Job{})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /jobs", getJobsHandler)
	mux.HandleFunc("POST /jobs", createJobHandler)

	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}


func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	var jobs []models.Job
	result := db.Order("created_at desc").Limit(50).Find(&jobs)
	if result.Error != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}


func createJobHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Job
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}


	if input.PostedAt.IsZero() {
		input.PostedAt = time.Now()
	}


	result := db.Create(&input)
	if result.Error != nil {
		log.Printf("[Skip] Duplicate found: %s", input.URL)
		http.Error(w, "Job already exists", http.StatusConflict) 
		return
	}

	log.Printf("[New] Saved job: %s", input.Title)

	go notifications.SendJobAlert(input.Title, input.Company, input.URL, input.Location)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "saved", "id": fmt.Sprint(input.ID)})
}