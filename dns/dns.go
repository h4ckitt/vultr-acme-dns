package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const vultrAPIBaseURL = "https://api.vultr.com/v2"

type DNSRecord struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
}

func CreateTXTRecord(apiKey, domain, name, value string) error {
	record := DNSRecord{
		Name: name,
		Type: "TXT",
		Data: value,
		TTL:  300,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	log.Println(string(payload))

	log.Printf("Creating DNS Record For Domain %s", domain)
	url := fmt.Sprintf("%s/domains/%s/records", vultrAPIBaseURL, domain)

	log.Printf("Calling URL: %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", body)
	}

	// Wait for DNS propagation
	time.Sleep(20 * time.Second)

	return nil
}

func DeleteTXTRecord(apiKey, domain, name string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/domains/%s/records", vultrAPIBaseURL, domain), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", body)
	}

	var records struct {
		Records []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"records"`
	}
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	var recordID string
	for _, record := range records.Records {
		if record.Name == "_acme-challenge" && record.Type == "TXT" {
			recordID = record.ID
			break
		}
	}

	if recordID == "" {
		return fmt.Errorf("TXT record not found for deletion")
	}

	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s/domains/%s/records/%s", vultrAPIBaseURL, domain, recordID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", body)
	}

	return nil
}
