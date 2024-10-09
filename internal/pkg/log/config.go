package log

// ConfigLogger main struct for config logger.
type ConfigLogger struct {
	// Level logger (DEBUG, INFO, WARN, ERROR, PANIC).
	Level string

	// EncodingType: console or json.
	EncodingType string

	// EnableCaller show field caller in logs.
	EnableCaller bool

	// EnableStacktrace show stacktrace.
	EnableStacktrace bool
}

// MergeDefault sets default value if no values are set.
func (c *ConfigLogger) MergeDefault() {
	if c.Level == "" {
		c.Level = "INFO"
	}
	if c.EncodingType == "" {
		c.EncodingType = "console"
	}
}
