package abclog

type Logger interface {
	With(keysAndValues ...any) Logger

	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	DPanic(msg string, keysAndValues ...any)
	Panic(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
}
