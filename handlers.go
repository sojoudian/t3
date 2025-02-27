package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type TimeResponse struct {
	TorontoTime    string `json:"toronto_time"`
	TehranTime     string `json:"tehran_time"`
	TorontoTimeStr string `json:"toronto_time_str"`
	TehranTimeStr  string `json:"tehran_time_str"`
}

type ConversionRequest struct {
	City   string `json:"city"`
	Hour   int    `json:"hour"`
	Minute int    `json:"minute"`
}

type ConversionResponse struct {
	SourceCity string `json:"source_city"`
	SourceTime string `json:"source_time"`
	TargetCity string `json:"target_city"`
	TargetTime string `json:"target_time"`
}

// GetCurrentTimeHandler returns the current time in Toronto and Tehran
func GetCurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for browser compatibility
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	// Load location (timezone) information
	torontoLoc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading Toronto timezone: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		http.Error(w, "Error loading Tehran timezone: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get current UTC time and convert to the respective time zones
	now := time.Now().UTC()
	torontoTime := now.In(torontoLoc)
	tehranTime := now.In(tehranLoc)
	
	response := TimeResponse{
		TorontoTime:    torontoTime.Format(time.RFC3339),
		TehranTime:     tehranTime.Format(time.RFC3339),
		TorontoTimeStr: torontoTime.Format("3:04 PM - January 2, 2006"),
		TehranTimeStr:  tehranTime.Format("3:04 PM - January 2, 2006"),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// ConvertTimeHandler converts a given time between Toronto and Tehran
func ConvertTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for browser compatibility
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse request body
	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Load location (timezone) information
	torontoLoc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading Toronto timezone: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		http.Error(w, "Error loading Tehran timezone: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get today's date
	now := time.Now().UTC()
	year, month, day := now.Date()
	
	var sourceTime time.Time
	var targetTime time.Time
	var sourceCity, targetCity string
	
	// Create a time.Time object for the selected hour and minute
	if req.City == "Toronto" {
		sourceCity = "Toronto"
		targetCity = "Tehran"
		// Create time in Toronto timezone
		sourceTime = time.Date(year, month, day, req.Hour, req.Minute, 0, 0, torontoLoc)
		// Convert to Tehran time
		targetTime = sourceTime.In(tehranLoc)
	} else if req.City == "Tehran" {
		sourceCity = "Tehran"
		targetCity = "Toronto"
		// Create time in Tehran timezone
		sourceTime = time.Date(year, month, day, req.Hour, req.Minute, 0, 0, tehranLoc)
		// Convert to Toronto time
		targetTime = sourceTime.In(torontoLoc)
	} else {
		http.Error(w, "Invalid city: must be 'Toronto' or 'Tehran'", http.StatusBadRequest)
		return
	}
	
	// Format the times in 12-hour format with AM/PM
	sourceTimeStr := sourceTime.Format("3:04 PM")
	targetTimeStr := targetTime.Format("3:04 PM")
	
	response := ConversionResponse{
		SourceCity: sourceCity,
		SourceTime: sourceTimeStr,
		TargetCity: targetCity,
		TargetTime: targetTimeStr,
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
