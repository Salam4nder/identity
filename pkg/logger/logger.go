package logger

//TODO

// Logger is a generic interface for logging libraries that
// support basic logging functionality.
type Logger interface {
	// Debug logs a message at debug level.
	Debug(msg string, fields ...Field)

	// Info logs a message at info level.
	Info(msg string, fields ...Field)

	// Warn logs a message at warn level.
	Warn(msg string, fields ...Field)

	// Error logs a message at error level.
	Error(msg string, fields ...Field)

	// WithField returns a new logger with the specified field.
	WithField(key string, value interface{}) Logger

	// WithFields returns a new logger with the specified fields.
	WithFields(fields Fields) Logger

	// WithError returns a new logger with the specified error.
	WithError(err error) Logger

	// Sync flushes any buffered log entries.
	Sync() error
}

// Fields is a map of key-value pairs that can be passed to a logger
// to include additional context in a log message.
type Fields map[string]interface{}

// Field is a key-value pair that can be passed to a logger
// to include additional context in a log message.
type Field struct {
	Key   string
	Value interface{}
}
