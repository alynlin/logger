package logger

import "context"

func Debug(args ...interface{}) {
	log.Debug(args...)
}
func Info(args ...interface{}) {
	log.Info(args...)
}
func Warn(args ...interface{}) {
	log.Warn(args...)
}
func Error(args ...interface{}) {
	log.Error(args...)
}
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}
func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}
func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}
func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}
func Fatalf(template string, args ...interface{}) {
	log.Fatalf(template, args...)
}

func InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	traceID, spanID := FromContext(ctx)
	log.With("trace_id", traceID, "span_id", spanID).Infof(msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...interface{}) {
	traceID, spanID := FromContext(ctx)
	log.With("trace_id", traceID, "span_id", spanID).Errorf(msg, args...)
}
