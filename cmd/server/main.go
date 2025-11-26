package main
import (
	"devlink/internal/database"
	"devlink/internal/email"
	"devlink/internal/api"
	"log"
	"net/http"
	"os"
)
func main() {
	apiSecret := os.Getenv("API_SECRET")
	if apiSecret == "" {
		log.Fatal("API_SECRET environment variable is required")
	}
	var err error
	db, err := database.New("./data/devlink.db")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	emailService, err := email.NewService(
		os.Getenv("AWS_REGION"),
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("SES_SENDER"),
	)
	if err != nil {
		log.Fatal(err)
	}

	server := &api.Server{
		Store: db,
		Email: emailService,
		APISecret: apiSecret,
	}

	// Routes
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("api/ingest", server.IngestJobsHandler)
	http.HandleFunc("api/subscribe", server.SubscribeHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting Devlink Backend on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
