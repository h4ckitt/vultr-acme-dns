package main

import (
	"fmt"
	"os"

	"vultr-dns/dns"
)

func main() {
	domain := os.Getenv("CERTBOT_DOMAIN")
	token := "2YOKMZXWLOT25NSDVEUFWKCCQIUOAWCXMZHA" //os.Getenv("VULTR_API_KEY")
	subdomain := fmt.Sprintf("_acme-challenge.%s", domain)

	err := dns.DeleteTXTRecord(token, domain, subdomain)
	if err != nil {
		fmt.Printf("Failed to delete TXT record: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("TXT record deleted successfully")
}
