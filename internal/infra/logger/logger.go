package logger

type Logger interface {
	D(format string, args ...interface{})
	I(format string, args ...interface{})
	W(format string, args ...interface{})
	E(format string, args ...interface{})
	C(format string, args ...interface{})
	Fields(fields Fields) Logger
	Err(err error) Logger
}

type Fields map[string]interface{}
