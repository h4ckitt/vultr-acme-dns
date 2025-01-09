package main

import (
	"log"
	"os"

	"vultr-dns/dns"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".certbotenv"); err != nil {
		log.Println("Failed To Find .certbotenv File, Please Create One With Contents VULTR_API_KEY=<your-vultr-api-key>")
		return
	}

	domain := os.Getenv("CERTBOT_DOMAIN")
	token := os.Getenv("VULTR_API_KEY")
	txtValue := os.Getenv("CERTBOT_VALIDATION")

	record, err := dns.FetchDNSRecord(domain, txtValue, token)
	if err != nil {
		return
	}

	if record.ID != "" {
		if err = dns.DeleteTXTRecord(domain, record.ID, token); err != nil {
			log.Printf("Failed To Delete TXT Record For %s: %v\n", domain, err)
		}
	} else {
		if err = dns.CreateTXTRecord(domain, "_acme-challenge", txtValue, token); err != nil {
			log.Printf("Failed To Create TXT Record For %s: %v\n", domain, err)
		}
	}
}
