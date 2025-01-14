package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Metric struct {
	MetricName string  `json:"metric_name"`
	Value      float64 `json:"value"`
}

// Payload structure
type MetricUpdatePayload struct {
	Program string   `json:"program"`
	Metrics []Metric `json:"metrics"`
}

// sends metrics to the Agent's aggregator
func SendMetrics(host string, port int, payload MetricUpdatePayload) error {

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %v", err)
	}

	// Create a new request
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/push", host, port), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	//fmt.Println("Metrics sent successfully")
	return nil
}
