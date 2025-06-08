package httpresponses

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
	Error *ErrorInfo `json:"error,omitempty"`
}

type ErrorInfo struct {
	Type   string `json:"type"`
	Details map[string]string `json:"details,omitempty"`
}

type Config struct {
	Logger *slog.Logger
}

var defaultConfig = Config{
	Logger: slog.Default(),
}

func SetConfig(cfg Config) {
	if cfg.Logger != nil {
		defaultConfig.Logger = cfg.Logger
	}
}

type StatusConfig struct {
	LogLevel slog.Level
	DefaultMessage string
	ErrorType string
}

var statusConfigMap = map[int] StatusConfig{
	// Success responses
	http.StatusOK: {
		DefaultMessage: "Request was successful",
		LogLevel: slog.LevelInfo,
	},
	http.StatusCreated: {
		DefaultMessage: "Resource created successfully",
		LogLevel: slog.LevelInfo,
	},
	http.StatusAccepted: {
		DefaultMessage: "Request accepted",
		LogLevel: slog.LevelInfo,
	},
	http.StatusNoContent: {
		DefaultMessage: "Request completed successfully",
		LogLevel: slog.LevelInfo,
	},

	// Client error responses
	http.StatusBadRequest: {
		DefaultMessage: "The request contains invalid data",
		LogLevel: slog.LevelWarn,
		ErrorType: "validation_error",
	},
	http.StatusUnauthorized: {
		DefaultMessage: "Authentication is required to access this resource",
		LogLevel: slog.LevelWarn,
		ErrorType: "authentication_error",
	},
	http.StatusForbidden: {
		DefaultMessage: "You do not have permission to access this resource",
		LogLevel: slog.LevelWarn,
		ErrorType: "authorization_error",
	},
	http.StatusNotFound: {
		DefaultMessage: "The requested resource was not found",
		LogLevel: slog.LevelInfo,
		ErrorType: "not_found",
	},
	http.StatusMethodNotAllowed: {
		DefaultMessage: "The requested method is not allowed for this resource",
		LogLevel: slog.LevelWarn,
		ErrorType: "method_not_allowed",
	},
	http.StatusConflict: {
		DefaultMessage: "The request could not be completed due to a conflict with the current state of the resource",
		LogLevel: slog.LevelWarn,
		ErrorType: "conflict",
	},
	http.StatusUnprocessableEntity: {
		DefaultMessage: "The request was well-formed but could not be processed due to semantic errors",
		LogLevel: slog.LevelWarn,
		ErrorType: "unprocessable_entity",
	},
	http.StatusTooManyRequests: {
		DefaultMessage: "Too many requests have been made in a given amount of time",
		LogLevel: slog.LevelWarn,
		ErrorType: "rate_limit_exceeded",
	},

	// Server error responses
	http.StatusInternalServerError: {
		DefaultMessage: "An unexpected error occurred on the server",
		LogLevel: slog.LevelError,
		ErrorType: "internal_server_error",
	},
	http.StatusNotImplemented: {
		DefaultMessage: "The requested functionality is not implemented",
		LogLevel: slog.LevelError,
		ErrorType: "not_implemented",
	},
	http.StatusBadGateway: {
		DefaultMessage: "The server received an invalid response from an upstream server",
		LogLevel: slog.LevelError,
		ErrorType: "bad_gateway",
	},
	http.StatusServiceUnavailable: {
		DefaultMessage: "The server is currently unable to handle the request due to temporary overload or maintenance",
		LogLevel: slog.LevelError,
		ErrorType: "service_unavailable",
	},
	http.StatusGatewayTimeout: {
		DefaultMessage: "The server did not receive a timely response from an upstream server",
		LogLevel: slog.LevelError,
		ErrorType: "gateway_timeout",
	},
	http.StatusHTTPVersionNotSupported: {
		DefaultMessage: "The server does not support the HTTP protocol version used in the request",
		LogLevel: slog.LevelError,
		ErrorType: "http_version_not_supported",
	},
	http.StatusVariantAlsoNegotiates: {
		DefaultMessage: "The server has an internal configuration error and cannot complete the request",
		LogLevel: slog.LevelError,
		ErrorType: "variant_also_negotiates",
	},
}

func AddStatusConfig(statusCode int, config StatusConfig) {
	if _, exists := statusConfigMap[statusCode]; !exists {
		statusConfigMap[statusCode] = config
	} else {
		defaultConfig.Logger.Warn("Status code already exists, updating configuration", slog.Int("statusCode", statusCode))
		statusConfigMap[statusCode] = config
	}
}

func getMessageForStatus(statusCode int, providedMessage string) string {
	if providedMessage != "" {
		return providedMessage
	}

	if config, exists := statusConfigMap[statusCode]; exists {
		return config.DefaultMessage
	}

	switch {
	case statusCode >= 200 && statusCode < 300:
		return "Request completed successfully"
	case statusCode >= 300 && statusCode < 400:
		return "Request requires further action"
	case statusCode >= 400 && statusCode < 500:
		return "Client error occurred"
	case statusCode >= 500:
		return "Server error occurred"
	default:
		return "Response completed"
	}
}

func validateStatusCode(statusCode int) int {
	if statusCode < 100 || statusCode > 599 {
		defaultConfig.Logger.Warn("Invalid HTTP status code provided, using 500",
			slog.Int("provided_code", statusCode))
		return http.StatusInternalServerError
	}
	return statusCode
}

func JSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, data interface{}, details map[string]string) {
	statusCode = validateStatusCode(statusCode)
	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	} else {
		ctx = context.Background()
	}
	var reqInfo map[string]interface{}
	message = getMessageForStatus(statusCode, message)

	if r != nil {
		ctx = r.Context()
		reqInfo = extractRequestInfo(r)
	} else {
		reqInfo = map[string]interface{}{
			"method":     "UNKNOWN",
			"path":       "UNKNOWN",
			"user_agent": "UNKNOWN",
			"remote_ip":  "UNKNOWN",
		}
		defaultConfig.Logger.Warn("JSON response called with nil request")
	}

	status := "success"
	var errorInfo *ErrorInfo

	if statusCode >= 400 {
		status = "error"
		
		// Get error type from configuration
		config, exists := statusConfigMap[statusCode]
		errorType := "unknown_error" // fallback
		if exists && config.ErrorType != "" {
			errorType = config.ErrorType
		}
		
		errorInfo = &ErrorInfo{
			Type:    errorType,
			Details: details,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	
	resp := Response{
		Status:  status,
		StatusCode:    statusCode,
		Message: message,
		Data:    data,
		Error:   errorInfo,
	}

	defaultConfig.Logger.Log(config.LogLevel, "HTTP response", slog.Int("statusCode", statusCode), slog.String("message", message))

	w.WriteHeader(statusCode)
	logAttrs := []slog.Attr{
		slog.Int("statusCode", statusCode),
		slog.String("status", status),
		slog.String("message", message),
		slog.String("method", reqInfo["method"].(string)),
		slog.String("path", reqInfo["path"].(string)),
		slog.String("user_agent", reqInfo["user_agent"].(string)),
		slog.String("remote_ip", reqInfo["remote_ip"].(string)),
	}

	if errorInfo != nil {
		logAttrs = append(logAttrs, 
			slog.String("error_type", errorInfo.Type),
			slog.Any("error_details", errorInfo.Details),
		)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		defaultConfig.Logger.ErrorContext(ctx, "Failed to encode JSON response",
			append(logAttrs, slog.Any("encoding_error", err))...)
		return
	}

	logLevel := slog.LevelInfo
	logMessage := "HTTP response sent"

	if config, exists := statusConfigMap[statusCode]; exists {
		logLevel = config.LogLevel
	} else {
		if statusCode >= 500 {
			logLevel = slog.LevelError
		} else if statusCode >= 400 {
			logLevel = slog.LevelWarn
		}
	}

	if statusCode >= 400 {
		if statusCode >= 500 {
			logMessage = "HTTP server error response sent"
		} else {
			logMessage = "HTTP client error response sent"
		}
	}

	defaultConfig.Logger.LogAttrs(ctx, logLevel, logMessage, logAttrs...)
}