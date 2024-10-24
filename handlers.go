package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func testSuitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error converting to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
