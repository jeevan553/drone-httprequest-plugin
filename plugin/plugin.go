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
	URL         string `envconfig:"PLUGIN_URL" required:"true"`
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

	// HTTP client setup with optional custom CA certificate
	client := &http.Client{Timeout: 10 * time.Second}
	tlsConfig := &tls.Config{}

	if cfg.CACertPath != "" {
		caCert, err := ioutil.ReadFile(cfg.CACertPath)
		if err != nil {
			fmt.Println("Error loading custom CA certificate:", err)
			return
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	// Prepare request body
	var requestBody *bytes.Reader
	if cfg.RequestBody != "" {
		requestBody = bytes.NewReader([]byte(cfg.RequestBody))
	} else {
		requestBody = bytes.NewReader(nil)
	}

	// Create HTTP request
	req, err := http.NewRequest(strings.ToUpper(cfg.HttpMethod), cfg.URL, requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set content type
	if cfg.ContentType != "" {
		req.Header.Set("Content-Type", cfg.ContentType)
	}

	// Set additional headers
	if cfg.Headers != "" {
		var headersMap map[string]string
		err := json.Unmarshal([]byte(cfg.Headers), &headersMap)
		if err != nil {
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
