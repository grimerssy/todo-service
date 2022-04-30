package logging

type Level byte

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Fields map[string]any

type Logger interface {
	Log(level Level, args ...any)
	Logln(level Level, args ...any)
	Logf(level Level, format string, args ...any)

	LogFields(level Level, fields Fields, args ...any)
	LogFieldsln(level Level, fields Fields, args ...any)
	LogFieldsf(level Level, fields Fields, format string, args ...any)
}
