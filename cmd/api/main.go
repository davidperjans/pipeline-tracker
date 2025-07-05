package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/davidperjans/pipeline-tracker/internal/middleware"
	"github.com/davidperjans/pipeline-tracker/internal/pipeline"
	"github.com/davidperjans/pipeline-tracker/internal/storage"
)

// Hanterar både GET och POST på /api/pipeline-runs
func pipelineRunsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetPipelineRuns(w, r)
	case http.MethodPost:
		handlePostPipelineRun(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostPipelineRun(w http.ResponseWriter, r *http.Request) {
	var run pipeline.PipelineRun
	if err := json.NewDecoder(r.Body).Decode(&run); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := pipeline.InsertPipelineRun(run); err != nil {
		http.Error(w, "Failed to insert run", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Pipeline run logged!"))
}

func handleGetPipelineRuns(w http.ResponseWriter, r *http.Request) {
	runs, err := pipeline.GetAllPipelineRuns()
	if err != nil {
		log.Printf("Error fetching runs: %v", err)
		http.Error(w, "Failed to fetch runs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runs)
}

func main() {
	// 1. Anslut till PostgreSQL
	if err := storage.ConnectToDB(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")

	// 2. Setup routes via ServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/api/pipeline-runs", pipelineRunsHandler)

	// 3. Wrap med middleware
	wrapped := middleware.RecoverPanic(middleware.RequestLogger(mux))

	// 4. Starta server med custom handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: wrapped,
	}

	go func() {
		log.Println("🚀 Server running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// 4. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
