package api

import (
	"devlink/internal/database"
	"devlink/internal/email"
	"devlink/internal/models"
	"encoding/json"
	"net/http"
)

type Server struct {
	Store     *database.Store
	Email     *email.Service
	APISecret string
}

func (s *Server) IngestJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := r.Header.Get("X-API-Token")
	if token != s.APISecret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req models.IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var newJobs []models.Job
	for _, job := range req.Jobs {
		inserted, err := s.Store.SaveJob(job)
		if err != nil && inserted {
			newJobs = append(newJobs, job)
		}
	}
	if len(newJobs) > 0 {
		go func() {
			subs, _ := s.Store.GetAllSubscribers()
			s.Email.Broadcast(subs, newJobs)
		}()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"new_jobs_count": len(newJobs),
	})
}

func (s *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var sub models.Subscriber
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if err := s.Store.SaveSubscriber(sub); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}