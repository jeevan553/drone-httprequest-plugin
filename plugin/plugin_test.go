package plugin

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestExecuteRequest(t *testing.T) {
	tests := []struct {
		name       string
		testType   string
		httpMethod string
		statusCode int
		response   string
	}{
		{
			name:       "Success Case",
			testType:   "success",
			httpMethod: "POST",
			statusCode: http.StatusCreated,
			response:   `{"message": "success"}`,
		},
		{
			name:       "Invalid Method",
			testType:   "invalid_method",
			httpMethod: "INVALID_METHOD",
			statusCode: http.StatusMethodNotAllowed,
			response:   `{"error": "method not allowed"}`,
		},
		{
			name:       "Server Error",
			testType:   "server_error",
			httpMethod: "POST",
			statusCode: http.StatusInternalServerError,
			response:   `{"error": "internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch tt.testType {
				case "success":
					if r.Method != tt.httpMethod {
						t.Errorf("Expected method %s, got %s", tt.httpMethod, r.Method)
					}
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))

				case "invalid_method":
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))

				case "server_error":
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(tt.response))

				default:
					t.Errorf("Unknown test type: %s", tt.testType)
				}
			}))
			defer mockServer.Close()

			// Set environment variables
			os.Setenv("PLUGIN_URL", mockServer.URL)
			os.Setenv("PLUGIN_HTTP_METHOD", tt.httpMethod)
			os.Setenv("PLUGIN_CONTENT_TYPE", "application/json")
			os.Setenv("PLUGIN_REQUEST_BODY", `{"title": "Hello", "body": "World", "userId": 1}`)

			// Execute the request
			ExecuteRequest()
		})
	}
}
