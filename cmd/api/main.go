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

	"github.com/davidperjans/pipeline-tracker/internal/pipeline"
	"github.com/davidperjans/pipeline-tracker/internal/storage"
)

func handlePipelineRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

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

func main() {
	// Connect to PostgreSQL database
	err := storage.ConnectToDb()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL")

	// Setup router
	http.HandleFunc("/api/pipeline-runs", handlePipelineRun)

	// Start HTTP server
	server := &http.Server{Addr: ":8080"}

	// Graceful shutdown handling
	go func() {
		log.Println("Server running at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: %v", err)
		}
	}()

	// Wait on Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
