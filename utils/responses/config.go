package responses

import "log/slog"

// Config holds configuration options for the httpresponses package.
type Config struct {
	Logger *slog.Logger
}

var defaultConfig = Config{
	Logger: slog.Default(),
}

// Only non-nil Logger will overwrite the default.
func SetConfig(cfg Config) {
	if cfg.Logger != nil {
		defaultConfig.Logger = cfg.Logger
	}
}
