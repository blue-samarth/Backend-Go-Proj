package responses

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func HTTPResponse(w http.ResponseWriter, r *http.Request, statusCode int, message string, data interface{}, details map[string]string) {
	statusCode = validateStatusCode(statusCode)

	var ctx context.Context
	if r != nil {
		ctx = r.Context()
	} else {
		ctx = context.Background()
	}

	var reqInfo RequestInfo
	if r != nil {
		reqInfo = extractRequestInfo(r)
	} else {
		defaultConfig.Logger.Warn("JSON response called with nil request")
	}

	message = getMessageForStatus(statusCode, message)

	status := "success"
	var errorInfo *ErrorInfo

	config, exists := statusConfigMap[statusCode]

	if statusCode >= 400 {
		status = "error"

		errorType := "unknown_error"
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
		Status:     status,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Error:      errorInfo,
	}

	// Determine log level
	logLevel := slog.LevelInfo
	if exists {
		logLevel = config.LogLevel
	} else if statusCode >= 500 {
		logLevel = slog.LevelError
	} else if statusCode >= 400 {
		logLevel = slog.LevelWarn
	}

	w.WriteHeader(statusCode)

	logAttrs := []slog.Attr{
		slog.Int("statusCode", statusCode),
		slog.String("status", status),
		slog.String("message", message),
		slog.String("method", reqInfo.Method),
		slog.String("path", reqInfo.Path),
		slog.String("user_agent", reqInfo.UserAgent),
		slog.String("remote_ip", reqInfo.RemoteIP),
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

	logMessage := "HTTP response sent"
	if statusCode >= 500 {
		logMessage = "HTTP server error response sent"
	} else if statusCode >= 400 {
		logMessage = "HTTP client error response sent"
	}

	defaultConfig.Logger.LogAttrs(ctx, logLevel, logMessage, logAttrs...)
}
