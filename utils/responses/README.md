# 📦 `responses` — A Structured HTTP Response Package for Go

This package standardizes JSON HTTP responses for Go APIs by providing unified response structures, configurable logging via `slog`, centralized HTTP status mappings, extracted request metadata, and one entry-point function for consistent API responses.

---

## 🏗️ `types.go` — Response Types

This file defines the core data structures used for standardized HTTP JSON responses in this package.

### 🔸 `Response`

Represents the standard structure for all HTTP JSON responses.

```go
type Response struct {
    Status     string      // "success" or "error"
    StatusCode int         // HTTP status code
    Message    string      // Human-readable message
    Data       interface{} // Optional payload data
    Error      *ErrorInfo  // Optional error details
}
```

**Status** indicates if the response is a `"success"` or `"error"`. 
- This field provides immediate clarity about the operation outcome without requiring clients to interpret HTTP status codes. 
- Think of this as the semantic meaning of your response, where a client can quickly determine if their request succeeded or failed without having to decode numerical status codes.

**StatusCode** contains the HTTP status code such as 200, 400, or 500. 
- While the Status field gives semantic meaning, this field enables programmatic handling based on standard HTTP conventions. 
- This dual approach allows both human-readable responses and machine-friendly processing, giving you the best of both worlds in API design.

**Message** provides a human-readable message for the client. 
- This field bridges the gap between technical status codes and user-friendly communication. 
- When you need to display something to an end user, this message can often be shown directly, while technical integrations can rely on the status code and error type for automated handling.

**Data** holds the main response payload and remains optional. 
- For error responses, this field typically stays null, while successful operations populate it with relevant information. 
- This separation allows clients to handle data and errors distinctly, making error handling more predictable and data processing cleaner.

**Error** contains detailed error information when applicable. 
- This field remains null for successful responses but provides structured error details when Status equals "error". 
- The structured approach helps clients understand not just that something went wrong, but specifically what went wrong and potentially how to fix it.

### 🔸 `ErrorInfo`

Provides structured details about an error.

```go
type ErrorInfo struct {
    Type    string            // Error type identifier (e.g., "validation_error")
    Details map[string]string // Additional error details (optional)
}
```

**Type** serves as a short string identifying the error category. 
- Common examples include "validation_error", "authentication_error", "rate_limit_error", and "internal_error". 
- This categorization enables clients to implement appropriate error handling strategies. 
- For instance, a client might retry after a delay for rate limit errors, redirect to login for authentication errors, or show field-specific messages for validation errors.

**Details** offers key-value pairs with extra error context. 
- This optional field allows you to provide specific information about what went wrong, such as which field failed validation or what constraint was violated. 
- This granular information proves invaluable for debugging and for providing helpful error messages to users.

### 🔸 `RequestInfo`

Holds extracted information from the HTTP request, useful for logging or tracing.

```go
type RequestInfo struct {
    Method    string // HTTP method (GET, POST, etc.)
    Path      string // Request path (URL.Path)
    UserAgent string // User-Agent header string
    RemoteIP  string // Client IP address
}
```

This structure captures essential request metadata that proves invaluable for debugging, monitoring, and security analysis. 
- The extracted information helps correlate responses with specific requests and provides context for troubleshooting issues in production environments. 
- When you encounter problems in production, having this information readily available in your logs can mean the difference between quick resolution and hours of investigation.

---

## ⚙️ `config.go` — Package Configuration

This file defines configuration options for the `responses` package, allowing you to customize logging behavior and integrate with your existing infrastructure.

### 🔸 `Config`

```go
type Config struct {
    Logger *slog.Logger
}
```

**Logger** accepts a pointer to a [`log/slog.Logger`](https://pkg.go.dev/log/slog#Logger) instance. 
- When set, this logger handles all logging within the package, enabling seamless integration with your application's logging strategy. 
- This design respects the principle that logging should be consistent across your entire application rather than each package maintaining its own logging approach.

### 🔸 `defaultConfig`

```go
var defaultConfig = Config{
    Logger: slog.Default(),
}
```

The package maintains a default configuration using Go's standard logger unless you explicitly override it. 
- This approach ensures the package works out of the box while remaining flexible for custom setups. 
- The principle here is convention over configuration - sensible defaults that work immediately, with the flexibility to customize when needed.

### 🔸 `SetConfig`

```go
func SetConfig(cfg Config)
```

Call this function to override the default logger configuration. 
- The function only replaces the default logger when you provide a non-nil `Logger` in the configuration struct.

**Example Usage:**
```go
import (
    "log/slog"
    "your-module-path/utils/responses"
)

// Create a custom logger with specific formatting
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
responses.SetConfig(responses.Config{Logger: logger})
```

This mechanism allows you to integrate the package with your application's preferred logging setup, whether you use structured JSON logging, specific output destinations, or custom log levels. 
- The beauty of this approach is that once configured, every response automatically uses your logging infrastructure without requiring additional setup in each handler.

---

## 🗺️ `status_map.go` — HTTP Status Code Mapping

This file centralizes the mapping of HTTP status codes to default messages, error types, and log levels. 
- This approach ensures that all API responses are consistent, maintainable, and easy to update as your application grows.

### ✅ Why Centralize Status Code Mapping?

**Consistency** emerges as the primary benefit of centralized status code mapping. 
- By defining all status code behaviors in one place, you ensure that every part of your application responds to clients in a uniform way. 
- A 400 Bad Request always carries the same default message and error type, regardless of which endpoint generates it. 
- This consistency builds trust with API consumers because they know what to expect.

**Maintainability** becomes significantly easier when you need to update messages, error types, or log levels for specific status codes. 
- Instead of hunting through your entire codebase for every place that returns a 401 Unauthorized, you make changes in a single location, and those changes propagate throughout your application automatically. 
- This centralization reduces the risk of inconsistencies and makes updates faster and more reliable.

**Extensibility** allows you to add support for new status codes or custom application codes without modifying existing response logic. 
- You simply add new entries to the mapping, and the existing infrastructure handles them appropriately. 
- This approach supports growth and evolution of your API without requiring architectural changes.

**Separation of Concerns** keeps your business logic clean by decoupling response formatting from core functionality. 
- Your handlers focus on processing requests and determining outcomes, while the response package handles the technical details of HTTP communication. 
- This separation makes your code more testable and easier to understand.

### 🔸 `StatusConfig`

```go
type StatusConfig struct {
    LogLevel       slog.Level
    DefaultMessage string
    ErrorType      string
}
```

**LogLevel** determines the severity level for logging this response. 
- Success responses typically use Info level, client errors use Warn level, and server errors use Error level. 
- This automatic log level assignment ensures consistent logging practices across your application. 
- When you review logs, you can quickly filter by severity to focus on the most important issues.

**DefaultMessage** provides the default message sent to clients when you don't specify a custom message. 
- These defaults save you from repeating common messages while still allowing customization when needed. 
- The default messages are carefully crafted to be informative but not overly technical, striking a balance between helpfulness and security.

**ErrorType** offers a short string categorizing the error, which proves useful for client-side error handling and analytics. 
- Common error types include "validation_error", "authentication_error", and "rate_limit_error". 
- These categories allow clients to implement sophisticated error handling strategies without parsing human-readable error messages.

### 🔸 `statusConfigMap`

The package includes a comprehensive map of HTTP status codes to their `StatusConfig`. 
- This mapping covers common success, client error, and server error codes, but you can extend it for custom codes specific to your application.

**Example Mapping:**
```go
var statusConfigMap = map[int]StatusConfig{
    http.StatusOK: {
        DefaultMessage: "Request was successful",
        LogLevel:       slog.LevelInfo,
    },
    http.StatusBadRequest: {
        DefaultMessage: "The request contains invalid data",
        LogLevel:       slog.LevelWarn,
        ErrorType:      "validation_error",
    },
    http.StatusUnauthorized: {
        DefaultMessage: "Authentication is required",
        LogLevel:       slog.LevelWarn,
        ErrorType:      "authentication_error",
    },
    // Additional mappings for comprehensive coverage...
}
```

### 🔸 Utility Functions

**`getMessageForStatus`** returns the appropriate message for a given status code. 
- The function uses your provided message if available, falls back to the mapped default message, or provides a generic fallback if neither exists. 
- This layered approach ensures that responses always include meaningful messages while respecting your customization preferences.

**`GetStatusConfig`** retrieves the `StatusConfig` for a given HTTP status code and returns a boolean indicating whether the configuration exists in the map. 
- This function allows you to check for supported status codes and handle unsupported ones appropriately, providing a safety net for edge cases.

---

## 🔍 `extract.go` — Request Information Extraction

This file provides helper functions to extract relevant information from incoming HTTP requests, which is essential for logging, tracing, and building informative API responses.

### ✅ Why These Functions Are Needed

**Accurate Client Identification** presents a significant challenge in modern web infrastructure. 
- Requests often pass through proxies, load balancers, content delivery networks, and API gateways before reaching your application. 
- The real client IP address may be hidden behind several layers of infrastructure. 
- Think of this like trying to identify who sent you a letter that passed through multiple postal sorting facilities - each facility might stamp their own address on the envelope, obscuring the original sender.

Accurately extracting the client's IP address proves crucial for several reasons:
- Security measures like rate limiting and audit trails
- Analytics for understanding user behavior
- Debugging issues that affect specific clients
- Without proper IP extraction, you might inadvertently rate limit your own load balancer or lose crucial debugging information when trying to trace problematic requests.

**Consistent Logging** ensures that your application generates uniform, structured logs that teams can easily search, filter, and analyze. 
- By extracting and structuring request metadata such as HTTP method, request path, user agent, and client IP, you create logs that provide meaningful context for monitoring and troubleshooting. 
- Consistent logging is like having a well-organized filing system - you can quickly find what you need when problems arise.

**Reusable Abstraction** prevents code duplication and reduces the risk of mistakes when extracting request information in multiple places throughout your application. 
- Centralizing this logic ensures that all parts of your application extract request information using the same reliable methods. 
- This approach follows the DRY (Don't Repeat Yourself) principle while ensuring accuracy and consistency.

### 🔸 `getClientIP(r *http.Request) string`

This function attempts to determine the real client IP address by checking headers in a specific order designed to handle common proxy configurations.

The function first examines the `X-Forwarded-For` header:
- This header may contain a comma-separated list of IP addresses
- When present, the first IP address usually represents the original client, while subsequent addresses represent intermediate proxies or load balancers
- This header is like a chain of custody document that tracks each step the request took to reach your server

If `X-Forwarded-For` is not available, the function checks the `X-Real-IP` header:
- This is another common header set by proxies and load balancers to indicate the original client IP address
- Different proxy software uses different header conventions, so checking multiple headers increases the chances of finding the true client IP

As a final fallback, the function uses `RemoteAddr` from the HTTP request:
- This represents the direct remote address from the TCP connection
- This address may include a port number and might represent a proxy rather than the actual client, but it provides a guaranteed fallback when other methods fail

If all extraction methods fail, the function returns the `RemoteAddr` value as-is:
- This ensures that you always receive some form of client identifier even in unusual network configurations
- This defensive approach ensures that your logging and security systems always have some client information to work with

This layered approach maximizes the chance of obtaining the true client IP address, even in complex network setups involving multiple layers of proxies, load balancers, and content delivery networks.

### 🔸 `extractRequestInfo(r *http.Request) RequestInfo`

This function builds a comprehensive `RequestInfo` struct containing essential request metadata. The extracted information includes:
- The HTTP method (GET, POST, PUT, DELETE, etc.)
- The request path from the URL
- The User-Agent string identifying the client software
- The client IP address obtained through the intelligent IP detection logic

This structured approach to request information extraction provides several benefits:
- The information proves useful for logging and tracing, allowing you to correlate responses with specific requests and understand user behavior patterns
- The metadata can be included in API responses for debugging purposes, helping developers understand how their requests are being processed
- The information also supports audit and security requirements by providing a complete record of who made what requests and when

---

## 📤 `httpresponses.go` — Main HTTP Response Handler

This file contains the core function for sending standardized JSON HTTP responses in your API. 
- It brings together all the abstractions from the package to ensure every response is consistent, informative, and well-logged.

### ✅ Why This Function Is Needed

**Consistency** emerges as the primary benefit of centralizing HTTP response formatting. 
- When every endpoint uses the same response function, your entire API returns responses in an identical structure. 
- This consistency makes your API predictable and easier to consume, whether by frontend applications, mobile clients, or third-party integrators. 
- Think of this like having a standard format for all business letters - recipients know exactly where to find the information they need.

**Error Handling** becomes systematic and comprehensive when you centralize response logic. 
- The function automatically sets appropriate error types and details for error responses, reducing the boilerplate code required in individual handlers and ensuring that error information is always complete and properly formatted. 
- This approach prevents the common problem of inconsistent error responses that can confuse clients and make debugging difficult.

**Logging** provides crucial observability for production systems. 
- The function gathers request and response metadata for structured logging, making debugging and monitoring significantly easier. 
- Every response gets logged with consistent information, enabling you to track system behavior, identify patterns, and troubleshoot issues effectively. 
- This comprehensive logging approach transforms your application from a black box into a transparent, observable system.

**Security and Best Practices** are automatically enforced through the centralized response function. 
- The function sets appropriate headers such as `Content-Type`, `X-Content-Type-Options`, and `Cache-Control` to prevent common vulnerabilities like MIME type sniffing and inadvertent caching of sensitive information. 
- These security measures are applied consistently without requiring developers to remember to set them in each handler.

**Extensibility** allows you to easily integrate with custom loggers, add new response metadata, or modify response behavior across your entire application by making changes in a single location. 
- This centralized approach makes it easy to evolve your API's response format or add new features without touching every handler.

### 🔸 `HTTPResponse`

```go
func HTTPResponse(
    w http.ResponseWriter,
    r *http.Request,
    statusCode int,
    message string,
    data interface{},
    details map[string]string,
)
```

**w** represents the HTTP response writer that will receive the formatted response.

**r** provides the HTTP request context used for extracting request metadata and correlation information.

**statusCode** specifies the HTTP status code to send, which gets mapped to appropriate default messages and error types.

**message** allows you to provide a custom message that overrides the default message for the status code. 
- Pass an empty string to use the default message.

**data** contains the payload to include in the response. 
- For successful operations, this typically contains the requested information or confirmation of the completed action. 
- For error responses, this usually remains nil.

**details** provides optional error details used specifically for error responses. 
- This map allows you to include specific information about what went wrong, such as which field failed validation or what constraint was violated.

### 🔸 Internal Function Flow

The function follows a systematic process to ensure every response meets quality and consistency standards.

**Validation and Mapping** begins the process:
- Validates the provided status code and maps it to the appropriate default message and error type
- This step ensures that unsupported status codes are handled gracefully and that every response includes appropriate metadata
- The function acts as a quality gate, ensuring that responses always meet your standards

**Request Information Extraction** captures essential request metadata:
- Includes HTTP method, request path, user agent, and client IP address
- This information becomes crucial for logging, debugging, and audit trails
- By capturing this information automatically, the function ensures that your logs always include the context needed for effective troubleshooting

**Response Structure Building** creates a complete `Response` struct:
- Populates all relevant fields according to the status code and provided parameters
- The function determines whether to mark the response as success or error based on the status code and populates error information accordingly
- This systematic approach ensures that responses are always complete and properly structured

**Security Header Setting** adds important HTTP headers to protect against common vulnerabilities:
- Sets `Content-Type: application/json` to ensure proper content type handling
- Adds `X-Content-Type-Options: nosniff` to prevent MIME type sniffing attacks
- Includes `Cache-Control: no-store` to prevent caching of potentially sensitive API responses
- These headers provide defense in depth against common web vulnerabilities

**HTTP Status Code Writing** sends the appropriate HTTP status code to the client before writing the response body:
- This step ensures that clients receive the correct status code for programmatic handling while the response body provides detailed information

**JSON Encoding and Response Writing** safely serializes the response structure to JSON and writes it to the response writer:
- The function handles encoding errors gracefully by logging them appropriately, ensuring that problems with response serialization don't go unnoticed

**Structured Logging** records the response with comprehensive structured attributes:
- Includes all request metadata, response status, and error details when present
- This logging provides the observability needed for production monitoring and debugging
- The structured format makes logs searchable and analyzable, supporting both human investigation and automated monitoring

**Error Handling** ensures that JSON encoding errors are logged at the error level:
- Provides visibility into serialization issues that might otherwise go unnoticed
- This defensive approach helps you identify and fix problems before they impact users

### 🔸 Example Usage

**Successful Response:**
```go
// Returning user profile data
responses.HTTPResponse(
    w, r,
    http.StatusOK,
    "User profile retrieved successfully",
    userProfile,
    nil,
)
```

**Error Response with Custom Message:**
```go
// Validation error with specific details
responses.HTTPResponse(
    w, r,
    http.StatusBadRequest,
    "User registration failed due to invalid input",
    nil,
    map[string]string{
        "field": "email",
        "reason": "invalid format",
    },
)
```

**Error Response with Default Message:**
```go
// Using default message for the status code
responses.HTTPResponse(
    w, r,
    http.StatusUnauthorized,
    "", // Empty string uses default message
    nil,
    nil,
)
```

This function serves as the single entry point for all HTTP responses in your API:
- Ensures every response is secure, consistent, and well-logged while reducing repetitive code in your handlers
- The centralized approach makes your API more maintainable and provides the foundation for consistent client experiences across your entire application

---

## 🎯 Summary

This package provides a comprehensive solution for building robust, production-ready Go APIs through standardized response handling. The architecture ensures:
- Uniform HTTP response formats across your entire application
- Secure response headers that protect against common vulnerabilities
- Meaningful and structured logs that support effective monitoring and debugging
- Reduced boilerplate code in route handlers
- Easy extension and integration with existing systems

The package embodies the principle of centralized response handling:
- Consistency, security, and observability are automatically enforced rather than being left to individual developers to implement correctly in each handler
- This approach significantly reduces the cognitive load on developers while improving the overall quality and maintainability of your API infrastructure

The modular design allows you to adopt the package incrementally:
- Start with basic response standardization and gradually leverage more advanced features like custom logging integration and detailed error categorization
- Whether you're building a simple API or a complex microservices architecture, this package provides the foundation for reliable, professional HTTP communication

Understanding this package architecture helps you build better APIs:
- Provides consistent, secure, and observable response handling that scales with your application's growth and complexity
- The investment in proper response handling pays dividends in reduced debugging time, improved client experience, and easier maintenance as your API evolves