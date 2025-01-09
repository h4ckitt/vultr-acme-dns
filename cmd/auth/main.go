package main

import (
	"fmt"
	"log"
	"os"

	"vultr-dns/dns"
)

func main() {
	domain := os.Getenv("CERTBOT_DOMAIN")
	token := "2YOKMZXWLOT25NSDVEUFWKCCQIUOAWCXMZHA" //os.Getenv("VULTR_API_KEY")
	txtValue := os.Getenv("CERTBOT_VALIDATION")
	subdomain := fmt.Sprintf("_acme-challenge.%s", domain)

	log.Println(token)

	err := dns.CreateTXTRecord(token, domain, subdomain, txtValue)
	if err != nil {
		fmt.Printf("Failed to create TXT record: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("TXT record created successfully")
}
