package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type APIResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

const (
	maxConcurrency = 50                // Matches instance concurrency
	timeout        = 120 * time.Second // Matches instance timeout
)

var (
	semaphore = make(chan struct{}, maxConcurrency)
	mu        sync.Mutex
)

func applyReferralHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MACUKK")
	// Acquire semaphore slot
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	// Simulate timeout handling
	done := make(chan struct{})
	go func() {
		defer close(done)
		fmt.Println("Starting processing...")
		// Artificial delay to simulate processing
		time.Sleep(2 * time.Second) // Adjust delay as needed
		fmt.Println("Processing completed...")

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// // Simulate occasional 504 Gateway Timeout
		// if time.Now().UnixNano()%20 == 0 { // Randomized failure
		// 	http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		// 	return
		// }

		response := APIResponse{
			Message: "Referral processed successfully",
			Status:  "success",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to write response:", err)
		}
	}()

	select {
	case <-done:
		// Completed successfully
		fmt.Println("Done")
	case <-time.After(timeout):
		http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		fmt.Println("ERROR")
	}
}

func main() {
	http.HandleFunc("/v1/referrals/apply", applyReferralHandler)
	log.Println("Mock API server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
