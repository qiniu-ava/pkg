package ava_grpc

import (
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

// assign log level for each grpc response code
var codeMappings = map[codes.Code]zapcore.Level{
	codes.OK:                 zapcore.InfoLevel,
	codes.Canceled:           zapcore.ErrorLevel,
	codes.Unknown:            zapcore.WarnLevel,
	codes.InvalidArgument:    zapcore.ErrorLevel,
	codes.DeadlineExceeded:   zapcore.DPanicLevel,
	codes.NotFound:           zapcore.ErrorLevel,
	codes.PermissionDenied:   zapcore.WarnLevel,
	codes.ResourceExhausted:  zapcore.DPanicLevel,
	codes.FailedPrecondition: zapcore.ErrorLevel,
	codes.Aborted:            zapcore.WarnLevel,
	codes.OutOfRange:         zapcore.WarnLevel,
	codes.Unimplemented:      zapcore.ErrorLevel,
	codes.Internal:           zapcore.ErrorLevel,
	codes.Unavailable:        zapcore.ErrorLevel,
	codes.DataLoss:           zapcore.DPanicLevel,
	codes.Unauthenticated:    zapcore.WarnLevel,
}

// CodeToLevel is going to be used with github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
func CodeToLevel(code codes.Code) zapcore.Level {
	return codeMappings[code]
}
