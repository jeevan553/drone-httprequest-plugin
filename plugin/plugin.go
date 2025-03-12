package plugin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config struct for environment variables
type Config struct {
	Url         string `envconfig:"PLUGIN_URL" required:"true"`
	HttpMethod  string `envconfig:"PLUGIN_HTTP_METHOD" required:"true"`
	Headers     string `envconfig:"PLUGIN_HEADERS"`
	ContentType string `envconfig:"PLUGIN_CONTENT_TYPE"`
	RequestBody string `envconfig:"PLUGIN_REQUEST_BODY"`
	CACertPath  string `envconfig:"PLUGIN_CA_CERT_PATH"`
}

// ExecuteRequest performs the HTTP request
func ExecuteRequest() {
	var cfg Config

	// Load environment variables
	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Println("Error parsing environment variables:", err)
		return
	}

	// Load CA certificates
	caCertPool, err := LoadCACertificates(cfg.CACertPath)
	if err != nil {
		fmt.Println("Error loading CA certificate:", err)
		return
	}

	// HTTP client setup
	tlsConfig := &tls.Config{}
	if caCertPool != nil {
		tlsConfig.RootCAs = caCertPool
	}
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
		Timeout:   10 * time.Second,
	}

	// Prepare request body
	requestBody := bytes.NewReader([]byte(cfg.RequestBody))

	// Create HTTP request
	req, err := http.NewRequest(strings.ToUpper(cfg.HttpMethod), cfg.Url, requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	if cfg.ContentType != "" {
		req.Header.Set("Content-Type", cfg.ContentType)
	}
	if cfg.Headers != "" {
		var headersMap map[string]string
		if err := json.Unmarshal([]byte(cfg.Headers), &headersMap); err != nil {
			fmt.Println("Invalid headers format:", err)
			return
		}
		for key, value := range headersMap {
			req.Header.Set(key, value)
		}
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Print response
	fmt.Println("Response Code:", resp.StatusCode)
	fmt.Println("Response Body:", string(body))
}

// LoadCACertificates loads system and custom CA certificates
func LoadCACertificates(caCertPath string) (*x509.CertPool, error) {
	// Load system CA pool
	sysCertPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Println("Warning: Failed to load system CA certificates:", err)
		sysCertPool = x509.NewCertPool()
	}

	// Load custom CA if provided
	if caCertPath == "" {
		return sysCertPool, nil
	}

	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read CA certificate file: %w", err)
	}

	if ok := sysCertPool.AppendCertsFromPEM(caCert); !ok {
		fmt.Println("Warning: Failed to append custom CA certificate")
	} else {
		fmt.Println("Custom CA certificate loaded successfully")
	}

	return sysCertPool, nil
}
