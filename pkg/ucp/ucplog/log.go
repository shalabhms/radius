/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ucplog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/radius-project/radius/pkg/version"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Radius uses the Zapr: https://github.com/go-logr/zapr which implements a logr interface
// for a zap log sink
const (
	DefaultLoggerName string = "radius"
	LogLevel          string = "RADIUS_LOGGING_LEVEL" // Env variable that determines the log level
	LogProfile        string = "RADIUS_LOGGING_JSON"  // Env variable that determines the logger config presets
)

// Log levels
const (
	// These log levels can be used with `logger.V(...)` to **add** verbosity to a log message.
	// These DO NOT map to the underlying Zap log levels.
	//
	// Please read the documentation of logger.V(...) before modifying these values.
	//
	// You CANNOT use these levels to make a log message more severe, which is why we don't
	// define log levels for warning or error. Use logger.Error() for errors.

	// LevelInfo is the default.
	LevelInfo int = 0

	// LevelDebug should be used for messages that should not be shown in production.
	LevelDebug int = 1
)

const (
	VerbosityLevelInfo  string = "INFO"
	VerbosityLevelDebug string = "DEBUG"
	VerbosityLevelError string = "ERROR"
	VerbosityLevelWarn  string = "WARN"
)

// Logger Profiles which determines the logger configuration
const (
	LoggerProfileProd    string = "production"
	LoggerProfileDev     string = "development"
	DefaultLoggerProfile        = LoggerProfileDev
)

func initLoggingConfig(options *LoggingOptions) (*zap.Logger, error) {
	var cfg zap.Config
	var loggerProfile, loggerLevel string

	// Define the logger profile and level based on the logger profile specified by RADIUS_LOGGING_JSON env variable or config files.
	// env variable takes precedence over config file settings.
	if options.Json {
		loggerProfile = LoggerProfileProd
	} else {
		loggerProfile = LoggerProfileDev
	}
	loggerProfileFromEnv := os.Getenv(LogProfile)
	if loggerProfileFromEnv != "" {
		loggerProfile = loggerProfileFromEnv
	}
	if strings.EqualFold(loggerProfile, LoggerProfileDev) {
		cfg = zap.NewDevelopmentConfig()
	} else if strings.EqualFold(loggerProfile, LoggerProfileProd) {
		cfg = zap.NewProductionConfig()
	} else {
		return nil, fmt.Errorf("invalid Radius Logger Profile set. Valid options are: %s, %s", LoggerProfileDev, LoggerProfileProd)
	}

	// Modify the default log level intialized by the profile preset if a custom value
	// is specified by config file or the "RADIUS_LOGGING_LEVEL" env variable. env variable takes precedence over config file settings.
	var logLevel int
	loggerLevel = options.Level
	logLevelFromEnv := os.Getenv(LogLevel)
	if logLevelFromEnv != "" {
		loggerLevel = logLevelFromEnv
	}

	if loggerLevel != "" {
		if strings.EqualFold(VerbosityLevelDebug, loggerLevel) {
			logLevel = int(zapcore.DebugLevel)
		} else if strings.EqualFold(VerbosityLevelInfo, loggerLevel) {
			logLevel = int(zapcore.InfoLevel)
		} else if strings.EqualFold(VerbosityLevelWarn, loggerLevel) {
			logLevel = int(zapcore.WarnLevel)
		} else if strings.EqualFold(VerbosityLevelError, loggerLevel) {
			logLevel = int(zapcore.ErrorLevel)
		} else {
			return nil, fmt.Errorf("invalid Radius Logger Level set. Valid options are: %s, %s, %s, %s", VerbosityLevelError, VerbosityLevelWarn, VerbosityLevelInfo, VerbosityLevelDebug)
		}
		cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(logLevel))
	}

	cfg.EncoderConfig.NameKey = "name"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "severity"
	cfg.EncoderConfig.TimeKey = "timestamp"

	// Build the logger config based on profile and custom presets
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize zap logger: %v", err)
	}

	return logger, nil
}

// NewLogger creates a new logger with zap logger implementation, with the given name and logging options,
// and returns a function to flush the logs before the server exits.
func NewLogger(name string, options *LoggingOptions) (logr.Logger, func(), error) {
	if name == "" {
		name = DefaultLoggerName
	}

	zapLogger, err := initLoggingConfig(options)
	if err != nil {
		return logr.Discard(), nil, err
	}
	logger := zapr.NewLogger(zapLogger).WithName(name)

	// Add the default resource key values, such as version, to new logger.
	logger = logger.WithValues(NewResourceObject(name)...)

	// The underlying zap logger needs to be flushed before server exits
	flushLogs := func() {
		err := zapLogger.Sync()
		if err != nil {
			logger.Error(err, "Unable to flush logs...")
		}
	}
	return logger, flushLogs, nil
}

// WrapLogContext adds key-value pairs to the context's logger for logging purposes.
func WrapLogContext(ctx context.Context, keyValues ...any) context.Context {
	logger := logr.FromContextOrDiscard(ctx)
	return logr.NewContext(ctx, logger.WithValues(keyValues...))
}

// Unwrap attempts to extract the underlying zap.Logger from a logr.Logger, returning nil if it fails.
func Unwrap(logger logr.Logger) *zap.Logger {
	underlier, ok := logger.GetSink().(zapr.Underlier)
	if ok {
		return underlier.GetUnderlying()
	}

	return nil
}

// FromContextOrDiscard returns a logger with trace and span IDs populated from the context if they exist.
// In order to get logger without span, use logr.FromContextOrDiscard(ctx context.Context).
func FromContextOrDiscard(ctx context.Context) logr.Logger {
	logger := logr.FromContextOrDiscard(ctx)

	// Populate trace id and span id when caller gets logger from context
	// because span id can be changed.
	sc := trace.SpanFromContext(ctx)
	if sc.SpanContext().HasTraceID() && sc.SpanContext().HasSpanID() {
		logger = logger.WithValues(
			LogFieldTraceId, sc.SpanContext().TraceID().String(),
			LogFieldSpanId, sc.SpanContext().SpanID().String(),
		)
	}
	return logger
}

// This function creates a new resource object with the given service name, hostname and version.
func NewResourceObject(serviceName string) []any {
	host, _ := os.Hostname()
	return []any{
		LogFieldServiceName, serviceName,
		LogFieldVersion, version.Channel(),
		LogFieldHostname, host,
	}
}
