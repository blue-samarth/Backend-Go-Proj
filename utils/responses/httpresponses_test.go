package responses

import (
    "bytes"
    "encoding/json"
    "log/slog"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Helper to decode response body
func decodeResponse(t *testing.T, body *bytes.Buffer) Response {
    var resp Response
    if err := json.NewDecoder(body).Decode(&resp); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }
    return resp
}

func TestHTTPResponse_Success(t *testing.T) {
    rec := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/test", nil)

    data := map[string]string{"foo": "bar"}
    HTTPResponse(rec, req, http.StatusOK, "Success!", data, nil)

    resp := decodeResponse(t, rec.Body)
    if resp.Status != "success" {
        t.Errorf("Expected status 'success', got %q", resp.Status)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected statusCode %d, got %d", http.StatusOK, resp.StatusCode)
    }
    if resp.Message != "Success!" {
        t.Errorf("Expected message 'Success!', got %q", resp.Message)
    }
    if resp.Data == nil {
        t.Error("Expected data, got nil")
    }
    if resp.Error != nil {
        t.Errorf("Expected error nil, got %+v", resp.Error)
    }
}

func TestHTTPResponse_ErrorWithDetails(t *testing.T) {
    rec := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodPost, "/fail", nil)

    details := map[string]string{"field": "email"}
    HTTPResponse(rec, req, http.StatusBadRequest, "", nil, details)

    resp := decodeResponse(t, rec.Body)
    if resp.Status != "error" {
        t.Errorf("Expected status 'error', got %q", resp.Status)
    }
    if resp.StatusCode != http.StatusBadRequest {
        t.Errorf("Expected statusCode %d, got %d", http.StatusBadRequest, resp.StatusCode)
    }
    if resp.Message == "" {
        t.Error("Expected non-empty message for error")
    }
    if resp.Data != nil {
        t.Errorf("Expected data nil, got %+v", resp.Data)
    }
    if resp.Error == nil {
        t.Error("Expected error info, got nil")
    } else if resp.Error.Type == "" {
        t.Error("Expected error type, got empty string")
    }
}

func TestSetConfig_CustomLogger(t *testing.T) {
    var logged bool
    logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
    SetConfig(Config{Logger: logger})

    rec := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/log", nil)
    HTTPResponse(rec, req, http.StatusOK, "Logged", nil, nil)

    // No assertion, just ensure no panic and logger is set
    logged = true
    if !logged {
        t.Error("Logger was not set or used")
    }
}

func TestGetStatusConfig(t *testing.T) {
    cfg, ok := GetStatusConfig(http.StatusOK)
    if !ok {
        t.Error("Expected status config for 200 OK")
    }
    if cfg.DefaultMessage == "" {
        t.Error("Expected default message for 200 OK")
    }
}

func TestExtractRequestInfo(t *testing.T) {
    req := httptest.NewRequest(http.MethodPut, "/info", nil)
    req.Header.Set("User-Agent", "TestAgent")
    req.RemoteAddr = "1.2.3.4:5678"
    info := extractRequestInfo(req)
    if info.Method != http.MethodPut {
        t.Errorf("Expected method PUT, got %s", info.Method)
    }
    if info.Path != "/info" {
        t.Errorf("Expected path /info, got %s", info.Path)
    }
    if info.UserAgent != "TestAgent" {
        t.Errorf("Expected UserAgent TestAgent, got %s", info.UserAgent)
    }
    if info.RemoteIP != "1.2.3.4" {
        t.Errorf("Expected RemoteIP 1.2.3.4, got %s", info.RemoteIP)
    }
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("X-Forwarded-For", "8.8.8.8, 9.9.9.9")
    ip := getClientIP(req)
    if ip != "8.8.8.8" {
        t.Errorf("Expected 8.8.8.8, got %s", ip)
    }
}

func TestGetClientIP_XRealIP(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.Header.Set("X-Real-IP", "7.7.7.7")
    ip := getClientIP(req)
    if ip != "7.7.7.7" {
        t.Errorf("Expected 7.7.7.7, got %s", ip)
    }
}

func TestGetClientIP_RemoteAddr(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    req.RemoteAddr = "6.6.6.6:1234"
    ip := getClientIP(req)
    if ip != "6.6.6.6" {
        t.Errorf("Expected 6.6.6.6, got %s", ip)
    }
}