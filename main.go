package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func flightAlertHandler(w http.ResponseWriter, r *http.Request) {
  var alert map[string]interface{}

  body, err := io.ReadAll(r.Body)

  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  err = json.Unmarshal(body, &alert)
  if err != nil {
    http.Error(w, "Invalid JSON", http.StatusBadRequest)
    return
  }

  // Add a timestamp field to the alert
	alert["receivedAt"] = time.Now().Format(time.RFC3339)

  	// Open the file in append mode, create if it doesn't exist
	file, err := os.OpenFile("alerts.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v\n", err)
		http.Error(w, "Unable to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

  // Convert the alert to pretty-printed JSON
	alertBytes, err := json.MarshalIndent(alert, "", "  ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v\n", err)
		http.Error(w, "Error processing alert", http.StatusInternalServerError)
		return
	}

	// Write the alert JSON to the file, appending a newline for clarity
	if _, err := file.Write(append(alertBytes, '\n')); err != nil {
		log.Printf("Error writing to file: %v\n", err)
		http.Error(w, "Error writing alert to file", http.StatusInternalServerError)
		return
	}

	// Log the alert to the console (optional)
	log.Printf("Received alert: %v\n", string(alertBytes))


  log.Printf("Received alert: %v\n", alert)

  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, "Received alert: %v\n", alert)
}

func main() {
  http.HandleFunc("/flight-alert", flightAlertHandler)

  log.Println("Starting server on :3001")
  log.Fatal(http.ListenAndServe(":3001", nil))
}
