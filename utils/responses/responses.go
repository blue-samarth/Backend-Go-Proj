package responses

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Response interface {
	Status string `json:"status"`
	StatusCode int    `json:"status_code"`
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