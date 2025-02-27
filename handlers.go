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
	torontoLoc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading Toronto timezone", http.StatusInternalServerError)
		return
	}
	
	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		http.Error(w, "Error loading Tehran timezone", http.StatusInternalServerError)
		return
	}
	
	now := time.Now()
	torontoTime := now.In(torontoLoc)
	tehranTime := now.In(tehranLoc)
	
	response := TimeResponse{
		TorontoTime:    torontoTime.Format(time.RFC3339),
		TehranTime:     tehranTime.Format(time.RFC3339),
		TorontoTimeStr: torontoTime.Format("3:04 PM - January 2, 2006"),
		TehranTimeStr:  tehranTime.Format("3:04 PM - January 2, 2006"),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConvertTimeHandler converts a given time between Toronto and Tehran
func ConvertTimeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req ConversionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	torontoLoc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading Toronto timezone", http.StatusInternalServerError)
		return
	}
	
	tehranLoc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		http.Error(w, "Error loading Tehran timezone", http.StatusInternalServerError)
		return
	}
	
	now := time.Now()
	var sourceTime time.Time
	var targetTime time.Time
	var sourceCity, targetCity string
	
	// Create a time.Time object for the selected hour
	if req.City == "Toronto" {
		sourceCity = "Toronto"
		targetCity = "Tehran"
		sourceTime = time.Date(now.Year(), now.Month(), now.Day(), req.Hour, req.Minute, 0, 0, torontoLoc)
		targetTime = sourceTime.In(tehranLoc)
	} else if req.City == "Tehran" {
		sourceCity = "Tehran"
		targetCity = "Toronto"
		sourceTime = time.Date(now.Year(), now.Month(), now.Day(), req.Hour, req.Minute, 0, 0, tehranLoc)
		targetTime = sourceTime.In(torontoLoc)
	} else {
		http.Error(w, "Invalid city", http.StatusBadRequest)
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
	json.NewEncoder(w).Encode(response)
}
