package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/syumai/workers"
)

// healthCheck calls the health check endpoint
func healthCheck(serverURL string) error {
	resp, err := http.Get(serverURL + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned non-OK status: %s", resp.Status)
	}

	return nil
}

// healthHandler handles the health check endpoint
func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Register health check endpoint
	http.HandleFunc("/health", healthHandler)

	// Default endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		msg := "Hello, Goooooo!"
		w.Write([]byte(msg))
	})

	// Initialize cron
	c := cron.New()

	// Schedule health check to run every minute
	_, err := c.AddFunc("@every 2s", func() {
		if err := healthCheck("http://localhost:9900"); err != nil {
			log.Printf("Error in health check: %v\n", err)
		} else {
			log.Printf("Health check successful at %v\n", time.Now().Format(time.RFC3339))
		}
	})

	if err != nil {
		panic("Failed to schedule health check: " + err.Error())
	}

	// Start cron
	c.Start()

	// Start the server
	workers.Serve(nil) // use http.DefaultServeMux
}
