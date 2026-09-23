package cmd

import "testing"

func TestExtractDomainsFromCSVUsesHeader(t *testing.T) {
	records := [][]string{{"IP", "ORIGIN", "TLS", "CERT_DOMAIN", "ALPN"}, {"1.2.3.4", "example", "TLS 1.3", "example.com", "h2"}}
	domains := extractDomainsFromCSV(records)
	if len(domains) != 1 || domains[0] != "example.com" {
		t.Fatalf("expected example.com, got %#v", domains)
	}
}

func TestExtractDomainsFromCSVRepairsUnquotedComma(t *testing.T) {
	records := [][]string{{"IP", "ORIGIN", "CERT_DOMAIN", "ALPN"}, {"1.2.3.4", "example", "foo.example", "bar.example", "h2"}}
	domains := extractDomainsFromCSV(records)
	if len(domains) != 1 || domains[0] != "foo.example,bar.example" {
		t.Fatalf("expected repaired domain, got %#v", domains)
	}
}

func TestExtractDomainsFromCSVIgnoresMissingHeader(t *testing.T) {
	if domains := extractDomainsFromCSV([][]string{{"IP", "TLS"}, {"1.2.3.4", "TLS 1.3"}}); len(domains) != 0 {
		t.Fatalf("expected no domains, got %#v", domains)
	}
}
