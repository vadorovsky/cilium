// Copyright 2016-2017 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"bufio"
	"bytes"
	"fmt"
	"log/syslog"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/cilium/cilium/pkg/logging/logfields"
	"github.com/cilium/cilium/pkg/metrics"

	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"sync/atomic"
)

const (
	fAddr  = "fluentd.address"
	fTag   = "fluentd.tag"
	fLevel = "fluentd.level"

	SLevel = "syslog.level"

	lAddr     = "logstash.address"
	lLevel    = "logstash.level"
	lProtocol = "logstash.protocol"

	Syslog   = "syslog"
	Fluentd  = "fluentd"
	Logstash = "logstash"
)

var (
	// DefaultLogger is the base logrus logger. It is different from the logrus
	// default to avoid external dependencies from writing out unexpectedly
	DefaultLogger = NewLogger()

	// DefaultLogLevel is the alternative we provide to Debug
	DefaultLogLevel = logrus.InfoLevel

	// syslogOpts is the set of supported options for syslog configuration.
	syslogOpts = map[string]bool{
		"syslog.level": true,
	}

	// fluentDOpts is the set of supported options for fluentD configuration.
	fluentDOpts = map[string]bool{
		fAddr:  true,
		fTag:   true,
		fLevel: true,
	}

	// logstashOpts is the set of supported options for logstash configuration.
	logstashOpts = map[string]bool{
		lAddr:     true,
		lLevel:    true,
		lProtocol: true,
	}

	// syslogLevelMap maps logrus.Level values to syslog.Priority levels.
	syslogLevelMap = map[logrus.Level]syslog.Priority{
		logrus.PanicLevel: syslog.LOG_ALERT,
		logrus.FatalLevel: syslog.LOG_CRIT,
		logrus.ErrorLevel: syslog.LOG_ERR,
		logrus.WarnLevel:  syslog.LOG_WARNING,
		logrus.InfoLevel:  syslog.LOG_INFO,
		logrus.DebugLevel: syslog.LOG_DEBUG,
	}
)

// Entry is a wrapper for the Logrus Entry. It exposes its public methods, but
// also adds a functionality of storing errors and warnings as Prometheus
// metrics. For detauls ebout Entry in Logrus see:
// https://godoc.org/github.com/sirupsen/logrus#Entry
type Entry struct {
	*logrus.Entry
	Logger    *Logger
	subsystem string
}

// WithError adds an error as single field (using the key defined in ErrorKey)
// to the Entry.
func (e *Entry) WithError(err error) *Entry {
	wrappedEntry := e.Entry.WithError(err)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: e.Logger,
	}
}

// WithField adds a single field to the Entry.
func (e *Entry) WithField(key string, value interface{}) *Entry {
	wrappedEntry := e.Entry.WithField(key, value)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: e.Logger,
	}
}

// WithFields adds a map of fields to the Entry.
func (e *Entry) WithFields(fields logrus.Fields) *Entry {
	wrappedEntry := e.Entry.WithFields(fields)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: e.Logger,
	}
}

// Warn logs a message at level Warn on the standard logger.
func (e *Entry) Warn(args ...interface{}) {
	e.Entry.Warn(args...)
	metrics.IncWarningsMetric(e.Logger.subsystem)
}

// Warning logs a message at level Warn on the standard logger.
func (e *Entry) Warning(args ...interface{}) {
	e.Entry.Warning(args...)
	metrics.IncWarningsMetric(e.Logger.subsystem)
}

// Error logs a message at level Error on the standard logger.
func (e *Entry) Error(args ...interface{}) {
	e.Entry.Error(args...)
	metrics.IncErrorsMetric(e.Logger.subsystem)
}

// Warnf logs a message at level Warn on the standard logger.
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.Warnf(format, args...)
	metrics.IncWarningsMetric(e.subsystem)
}

// Warningf logs a message at level Warn on the standard logger.
func (e *Entry) Warningf(format string, args ...interface{}) {
	e.Entry.Warningf(format, args)
	metrics.IncWarningsMetric(e.subsystem)
}

// Errorf logs a message at level Error on the standard logger.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.Errorf(format, args...)
	metrics.IncErrorsMetric(e.subsystem)
}

// NewEntry returns a new instance of Entry.
func NewEntry(logger *Logger) *Entry {
	wrappedEntry := logrus.NewEntry(logger.Logger)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: logger,
	}
}

// Logger is a wrapper for the Logrus Logger. It exposes its public methods, but
// also adds a functionality of storing errors and warnings as Prometheus
// metrics. For details about Logger in Logrus see:
// https://godoc.org/github.com/sirupsen/logrus#Logger
type Logger struct {
	*logrus.Logger
	subsystem string
}

// WithField adds a single field to the log entry.
func (l *Logger) WithField(key string, value interface{}) *Entry {
	wrappedEntry := l.Logger.WithField(key, value)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: l,
	}
}

// WithFields adds a map of fields to the log entry.
func (l *Logger) WithFields(fields logrus.Fields) *Entry {
	wrappedEntry := l.Logger.WithFields(fields)
	return &Entry{
		Entry:  wrappedEntry,
		Logger: l,
	}
}

// WithSubsystem adds an information about subsystem to the log entry.
func (l *Logger) WithSubsystem(subsystem string) *Entry {
	entry := l.WithField(logfields.LogSubsys, subsystem)
	entry.subsystem = subsystem
	return entry
}

// NewLogger returns a new instance of Logger with default attributes.
func NewLogger() *Logger {
	wrappedLogger := logrus.New()
	wrappedLogger.Formatter = setupFormatter()
	return &Logger{
		Logger: wrappedLogger,
	}
}

// setFireLevels returns a slice of logrus.Level values higher in priority
// and including level, excluding any levels lower in priority.
func setFireLevels(level logrus.Level) []logrus.Level {
	switch level {
	case logrus.PanicLevel:
		return []logrus.Level{logrus.PanicLevel}
	case logrus.FatalLevel:
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	case logrus.ErrorLevel:
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
	case logrus.WarnLevel:
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel}
	case logrus.InfoLevel:
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}
	case logrus.DebugLevel:
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel, logrus.DebugLevel}
	default:
		logrus.Infof("logrus level %v is not supported at this time; defaulting to info level", level)
		return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.InfoLevel}
	}
}

// SetupLogging sets up each logging service provided in loggers and configures
// each logger with the provided logOpts.
func SetupLogging(loggers []string, logOpts map[string]string, tag string, debug bool) error {
	// Set default logger to output to stdout if no loggers are provided.
	if len(loggers) == 0 {
		// TODO: switch to a per-logger version when we upgrade to logrus >1.0.3
		logrus.SetOutput(os.Stdout)
	}

	SetLogLevel(DefaultLogLevel)
	ToggleDebugLogs(debug)

	// always suppress the default logger so libraries don't print things
	logrus.SetLevel(logrus.PanicLevel)

	// Iterate through all provided loggers and configure them according
	// to user-provided settings.
	for _, logger := range loggers {
		valuesToValidate := getLogDriverConfig(logger, logOpts)
		switch logger {
		case Syslog:
			valuesToValidate := getLogDriverConfig(Syslog, logOpts)
			err := validateOpts(Syslog, valuesToValidate, syslogOpts)
			if err != nil {
				return err
			}
			setupSyslog(valuesToValidate, tag, debug)
		case Fluentd:
			err := validateOpts(logger, valuesToValidate, fluentDOpts)
			if err != nil {
				return err
			}
			setupFluentD(valuesToValidate, debug)
			//TODO - need to finish logstash integration.
		/*case Logstash:
		fmt.Printf("SetupLogging: in logstash case\n")
		err := validateOpts(logger, valuesToValidate, logstashOpts)
		fmt.Printf("SetupLogging: validating options for logstash complete\n")
		if err != nil {
			fmt.Printf("SetupLogging: error validating logstash opts %v\n", err.Error())
			return err
		}
		fmt.Printf("SetupLogging: about to setup logstash\n")
		setupLogstash(valuesToValidate)
		*/
		default:
			return fmt.Errorf("provided log driver %q is not a supported log driver", logger)
		}
	}

	return nil
}

// SetLogLevel sets the log level on DefaultLogger. This logger is, by
// convention, the base logger for package specific ones thus setting the level
// here impacts the default logging behaviour.
// This function is thread-safe when logging, reading DefaultLogger.Level is
// not protected this way, however.
func SetLogLevel(level logrus.Level) {
	DefaultLogger.SetLevel(level)
}

// ToggleDebugLogs switches on or off debugging logs. It will select
// DefaultLogLevel when turning debug off.
// It is thread-safe.
func ToggleDebugLogs(debug bool) {
	if debug {
		SetLogLevel(logrus.DebugLevel)
	} else {
		SetLogLevel(DefaultLogLevel)
	}
}

// setupSyslog sets up and configures syslog with the provided options in
// logOpts. If some options are not provided, sensible defaults are used.
func setupSyslog(logOpts map[string]string, tag string, debug bool) {
	logLevel, ok := logOpts[SLevel]
	if !ok {
		if debug {
			logLevel = "debug"
		} else {
			logLevel = "info"
		}
	}

	//Validate provided log level.
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		DefaultLogger.Fatal(err)
	}

	DefaultLogger.SetLevel(level)
	// Create syslog hook.
	h, err := logrus_syslog.NewSyslogHook("", "", syslogLevelMap[level], tag)
	if err != nil {
		DefaultLogger.Fatal(err)
	}
	// TODO: switch to a per-logger version when we upgrade to logrus >1.0.3
	logrus.AddHook(h)
}

// setupFormatter sets up the text formatting for logs output by logrus.
func setupFormatter() logrus.Formatter {
	fileFormat := new(logrus.TextFormatter)
	fileFormat.DisableTimestamp = true
	fileFormat.DisableColors = true
	switch os.Getenv("INITSYSTEM") {
	case "SYSTEMD":
		fileFormat.FullTimestamp = true
	default:
		fileFormat.TimestampFormat = time.RFC3339
	}

	// TODO: switch to a per-logger version when we upgrade to logrus >1.0.3
	return fileFormat
}

// setupFluentD sets up and configures FluentD with the provided options in
// logOpts. If some options are not provided, sensible defaults are used.
func setupFluentD(logOpts map[string]string, debug bool) {
	//If no logging level set for fluentd, use debug value if it is set.
	// Logging level set for fluentd takes precedence over debug flag
	// fluent.level provided.
	logLevel, ok := logOpts[fLevel]
	if !ok {
		if debug {
			logLevel = "debug"
		} else {
			logLevel = "info"
		}
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		DefaultLogger.Fatal(err)
	}

	hostAndPort, ok := logOpts[fAddr]
	if !ok {
		hostAndPort = "localhost:24224"
	}

	host, strPort, err := net.SplitHostPort(hostAndPort)
	if err != nil {
		DefaultLogger.Fatal(err)
	}
	port, err := strconv.Atoi(strPort)
	if err != nil {
		DefaultLogger.Fatal(err)
	}

	h, err := logrus_fluent.New(host, port)
	if err != nil {
		DefaultLogger.Fatal(err)
	}

	tag, ok := logOpts[fTag]
	if ok {
		h.SetTag(tag)
	}

	// set custom fire level
	h.SetLevels(setFireLevels(level))
	// TODO: switch to a per-logger version when we upgrade to logrus >1.0.3
	logrus.AddHook(h)
}

// setupLogstash sets up and configures Logstash with the provided options in
// logOpts. If some options are not provided, sensible defaults are used.
/// FIXME GH-1578 - needs to be tested with a working logstash setup.
//func setupLogstash(logOpts map[string]string) {
//	hostAndPort, ok := logOpts[lAddr]
//	if !ok {
//		hostAndPort = "172.17.0.2:999"
//	}
//
//	protocol, ok := logOpts[lProtocol]
//	if !ok {
//		protocol = "tcp"
//	}
//
//	h, err := logrustash.NewHook(protocol, hostAndPort, "cilium")
//	if err != nil {
//		DefaultLogger.Fatal(err)
//	}
//
//	DefaultLogger.AddHook(h)
//}

// validateOpts iterates through all of the keys in logOpts, and errors out if
// the key in logOpts is not a key in supportedOpts.
func validateOpts(logDriver string, logOpts map[string]string, supportedOpts map[string]bool) error {
	for k := range logOpts {
		if !supportedOpts[k] {
			return fmt.Errorf("provided configuration value %q is not supported as a logging option for log driver %s", k, logDriver)
		}
	}
	return nil
}

// getLogDriverConfig returns a map containing the key-value pairs that start
// with string logDriver from map logOpts.
func getLogDriverConfig(logDriver string, logOpts map[string]string) map[string]string {
	keysToValidate := make(map[string]string)
	for k, v := range logOpts {
		ok, err := regexp.MatchString(logDriver+".*", k)
		if err != nil {
			DefaultLogger.Fatal(err)
		}
		if ok {
			keysToValidate[k] = v
		}
	}
	return keysToValidate
}

// MultiLine breaks a multi line text into individual log entries and calls the
// logging function to log each entry
func MultiLine(logFn func(args ...interface{}), output string) {
	scanner := bufio.NewScanner(bytes.NewReader([]byte(output)))
	for scanner.Scan() {
		logFn(scanner.Text())
	}
}

// CanLogAt returns whether a log message at the given level would be
// logged by the given logger.
func CanLogAt(logger *Logger, level logrus.Level) bool {
	return GetLevel(logger) >= level
}

// GetLevel returns the log level of the given logger.
func GetLevel(logger *Logger) logrus.Level {
	return logrus.Level(atomic.LoadUint32((*uint32)(&logger.Level)))
}
