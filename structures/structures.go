package structures

type FetchDNSRecord struct {
	Records []DNSRecord `json:"records"`
}

type DNSRecord struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
	TTL  int    `json:"ttl"`
}
