package logger

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}
