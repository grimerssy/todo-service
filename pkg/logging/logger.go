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

type Fields map[string]interface{}

type Logger interface {
	Log(level Level, args ...interface{})
	Logln(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})

	LogFields(level Level, fields Fields, args ...interface{})
	LogFieldsln(level Level, fields Fields, args ...interface{})
	LogFieldsf(level Level, fields Fields, format string, args ...interface{})
}
