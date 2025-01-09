package dns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"vultr-dns/structures"
)

const vultrAPIBaseURL = "https://api.vultr.com/v2"

func CreateTXTRecord(domain, name, value, token string) error {
	record := structures.DNSRecord{
		Name: name,
		Type: "TXT",
		Data: value,
		TTL:  300,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	url := fmt.Sprintf("%s/domains/%s/records", vultrAPIBaseURL, domain)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
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

func DeleteTXTRecord(domain, recordID, token string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/domains/%s/records/%s", vultrAPIBaseURL, domain, recordID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
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

func FetchDNSRecord(domain, txtValue, token string) (record structures.DNSRecord, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/domains/%s/records", vultrAPIBaseURL, domain), http.NoBody)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed To Fetch DNS Records: %v\n", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Failed To Fetch DNS Records, API Error: %s\n", string(body))
		return
	}

	var res structures.FetchDNSRecord

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		log.Printf("Failed To Fetch DNS Records: %v\n", err)
		return
	}

	for _, result := range res.Records {
		if result.Type == "TXT" && strings.Contains(result.Data, txtValue) {
			record = result
			return
		}
	}

	return
}
