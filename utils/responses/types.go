package responses

// Response represents a standard HTTP JSON response structure.
type Response struct {
	Status     string      `json:"status"`               // "success" or "error"
	StatusCode int         `json:"statusCode"`           // HTTP status code
	Message    string      `json:"message"`              // Human-readable message
	Data       interface{} `json:"data,omitempty"`       // Payload data, optional
	Error      *ErrorInfo  `json:"error,omitempty"`      // Error details, optional
}

// ErrorInfo provides structured details about an error.
type ErrorInfo struct {
	Type    string            `json:"type"`               // Error type identifier (e.g., "validation_error")
	Details map[string]string `json:"details,omitempty"`  // Additional error details, optional
}

// RequestInfo holds extracted info from the HTTP request for logging or tracing.
type RequestInfo struct {
	Method    string // HTTP method (GET, POST, etc.)
	Path      string // Request path (URL.Path)
	UserAgent string // User-Agent header string
	RemoteIP  string // Client IP address
}
