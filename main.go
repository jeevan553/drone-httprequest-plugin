package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/jeevan553/drone-httprequest-plugin/plugin" // ✅ Correct Import
)

var urlToTest = "https://www.google.com"

func main() {
	fmt.Println("Starting Custom CA Verification...")

	// Load custom CA certificate
	homeDir, _ := os.UserHomeDir()
	caCertPath := homeDir + "/vcert.crt"

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("Error loading CA certificate from:", caCertPath, err)
	} else {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		testTLSConnection(urlToTest, caCertPool, "Custom CA")
	}

	fmt.Println("Testing system CA...")

	// Load system CA certificates
	sysCertPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("Error loading system CA pool:", err)
	} else {
		testTLSConnection(urlToTest, sysCertPool, "System CA (/etc/ssl/certs/ca-certificates.crt)")
	}

	fmt.Println("Merging System CA and Custom CA...")

	// Merge System CA with Custom CA
	if caCert != nil && sysCertPool != nil {
		sysCertPool.AppendCertsFromPEM(caCert)
		testTLSConnection(urlToTest, sysCertPool, "Merged (System + Custom CA)")
	}

	// Execute HTTP Request
	plugin.ExecuteRequest()
}

// testTLSConnection tests the TLS handshake with a given CA
func testTLSConnection(url string, caPool *x509.CertPool, testName string) {
	fmt.Println("Testing:", testName)
	tlsConfig := &tls.Config{RootCAs: caPool}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Request failed:", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("Request successful! HTTP Status:", resp.Status)
	}
}
