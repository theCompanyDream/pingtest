package main

import (
	"encoding/json"
	"net/http"
	"log"
	"github.com/theCompanyDream/pingtest/commands"
)

type CommandRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func main() {
	// Create a new ServeMux to handle multiple routes
	mux := http.NewServeMux()

	// Register the handlers
	mux.HandleFunc("/ping", handlePing)            // Endpoint for ping
	mux.HandleFunc("/", handleGetSystemInfo) // Endpoint for system info

	// Set up the server with the mux
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Assign the mux to the server's handler
	}

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(server.ListenAndServe())
}


func handlePing(w http.ResponseWriter, req *http.Request) {
	host := req.URL.Query().Get("host")
	if host == "" {
		http.Error(w, "Host parameter not found in query string", http.StatusBadRequest)
		return
	}
	result, err := commands.Ping(host)
	if err != nil {
		http.Error(w, "Error Running Ping Request", http.StatusBadRequest)
		log.Printf("Error pinging server %v", err)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func handleGetSystemInfo(w http.ResponseWriter, req *http.Request) {
	result, err := commands.GetSystemInfo()
	if err != nil {
		http.Error(w, "Error parsing query parameters", http.StatusBadRequest)
		log.Printf("Error parsing query parameters: %v", err)
		return
	}
	json.NewEncoder(w).Encode(result)
}
