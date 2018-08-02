package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golang/glog"
)

// DebugEnabled specifies if this package print debug information
var DebugEnabled = false

// LoggerT is the interface required to log test specific errors
// golang's standard `*testing.T` and `ginkgo.GinkgoTInterface` both implements this interface
type LoggerT interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// Logger is a struct which will help to call CITF specific logging functions
type Logger struct {
	T LoggerT
}

// below four methods are responsible to maintain uniform appearance of error message in logs

// ErrorMessageFromInterfaces prepares the message with one error and other interfaces
func (logger Logger) ErrorMessageFromInterfaces(err error, a ...interface{}) string {
	return fmt.Sprint(append(a, ":", err)...)
}

// ErrorMessageFromFormatString prepares the message over format string `message` and interfaces representing each format specifier
func (logger Logger) ErrorMessageFromFormatString(err error, message string, a ...interface{}) string {
	return fmt.Sprintf(strings.TrimSpace(message)+": "+err.Error(), a...)
}

// WritefDebugMessage formats according to a format specifier and writes to w only when DebugEnabled is true.
// A newline is always appended. It returns the number of bytes written and any write error encountered.
func (logger Logger) WritefDebugMessage(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if DebugEnabled {
		return fmt.Fprintf(w, format+"\n", a...)
	}
	return
}

// WritelnDebugMessage formats using the default formats for its operands and writes to w only when DebugEnabled.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) WritelnDebugMessage(w io.Writer, a ...interface{}) (n int, err error) {
	if DebugEnabled {
		return fmt.Fprintln(w, a...)
	}
	return
}

// PrintfDebugMessage formats according to a format specifier and writes to standard output only when DebugEnabled is true.
// //  A newline is always appended. It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintfDebugMessage(format string, a ...interface{}) (n int, err error) {
	return logger.WritefDebugMessage(os.Stdout, format, a...)
}

// PrintlnDebugMessage formats using the default formats for its operands and writes to standard output only when DebugEnabled.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintlnDebugMessage(a ...interface{}) (n int, err error) {
	return logger.WritelnDebugMessage(os.Stdout, a...)
}

// Log logs the info on testing variable (using `logger.T.Log`) if `logger.T` is not nil
// otherwise it uses `logger.PrintfDebugMessage`
func (logger Logger) Log(a ...interface{}) (n int, err error) {
	if logger.T != nil {
		logger.T.Log(a...)
	} else {
		return logger.PrintlnDebugMessage(a...)
	}
	return
}

// Logf logs the info on testing variable (using `logger.T.Logf`) if `logger.T` is not nil
// otherwise it uses `logger.PrintfDebugMessage`
func (logger Logger) Logf(format string, a ...interface{}) (n int, err error) {
	if logger.T != nil {
		logger.T.Logf(format, a...)
	} else {
		return logger.PrintfDebugMessage(format, a...)
	}
	return
}

// LogfDebugMessage logs the info if DebugEnabled. it does so on testing variable (using `logger.T.Logf`) if `logger.T` is not nil
// otherwise it uses `logger.PrintfDebugMessage`
func (logger Logger) LogfDebugMessage(format string, a ...interface{}) (n int, err error) {
	if DebugEnabled {
		return logger.Logf(format, a...)
	}
	return
}

// LogDebugMessage logs the info if DebugEnabled. it does so on testing variable (using `logger.T.Logf`) if `logger.T` is not nil
// otherwise it uses `logger.PrintfDebugMessage`
func (logger Logger) LogDebugMessage(a ...interface{}) (n int, err error) {
	if DebugEnabled {
		return logger.Log(a...)
	}
	return
}

// LogError logs error using `glog.Error` only when err is not nil.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogError(err error, message string) {
	if err != nil {
		glog.Error(logger.ErrorMessageFromInterfaces(err, message))
	}
}

// LogNonError logs info using `glog.Info` only when err is nil.
func (logger Logger) LogNonError(err error, a ...interface{}) {
	if err == nil {
		glog.Info(a...)
	}
}

// LogErrorf formats according to a format specifier and logs error using `glog.Error` only when err is not nil.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogErrorf(err error, message string, a ...interface{}) {
	if err != nil {
		glog.Error(logger.ErrorMessageFromFormatString(err, message, a...))
	}
}

// LogNonErrorf logs info using `glog.Infof` only when err is nil.
// formatting is taken care by `glog.Infof`
func (logger Logger) LogNonErrorf(err error, message string, a ...interface{}) {
	if err == nil {
		glog.Infof(message, a...)
	}
}

// logFatal is a plain function which does straight forward task of preparing and logging the error and exit.
// it is meant to be used in other functions of this package which requires above two requirements.
func (logger Logger) logFatal(err error, a ...interface{}) {
	glog.Fatal(logger.ErrorMessageFromInterfaces(err, a...))
}

// logFatalf is a plain function which does straight forward task of preparing and logging the error and exit.
// it is meant to be used in other functions of this package which requires above two requirements.
func (logger Logger) logFatalf(err error, message string, a ...interface{}) {
	glog.Fatal(logger.ErrorMessageFromFormatString(err, message, a...))
}

// LogFatal logs error using `glog.Error` only when err is not nil.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogFatal(err error, a ...interface{}) {
	if err != nil {
		logger.logFatal(err, a...)
	}
}

// LogFatalf formats according to a format specifier and logs error using `glog.Error` only when err is not nil.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogFatalf(err error, message string, a ...interface{}) {
	if err != nil {
		logger.logFatalf(err, message, a...)
	}
}

// LogTestError logs test error only when err is not nil.
// It does so using `logger.T.Error` if `logger.T` is not nil, otherwise it uses `glog.Fatal` (BE AWARE)
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogTestError(err error, a ...interface{}) {
	if err != nil {
		if logger.T != nil {
			logger.T.Error(a...)
		} else {
			logger.logFatal(err, a...)
		}
	}
}

// LogTestNonError logs info using `logger.Log` but only when err is nil.
func (logger Logger) LogTestNonError(err error, a ...interface{}) {
	if err == nil {
		logger.Log(a...)
	}
}

// LogTestErrorf formats according to a format specifier and logs error only when err is not nil.
// It does so using `logger.T.Errorf` if `logger.T` is not nil, otherwise it uses `glog.Fatal` (BE AWARE)
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) LogTestErrorf(err error, message string, a ...interface{}) {
	if err != nil {
		if logger.T != nil {
			logger.T.Errorf(message, a...)
		} else {
			logger.logFatalf(err, message, a...)
		}
	}
}

// LogTestNonErrorf logs info using `Logf` only when err is nil.
// formatting is taken care by `glog.Infof` pr `logger.T.Logf` accordingly.
func (logger Logger) LogTestNonErrorf(err error, message string, a ...interface{}) {
	if err == nil {
		logger.Logf(message, a...)
	}
}

// PrintError writes error message to os.StdErr only when err is not nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) PrintError(err error, message string) (n int, errr error) {
	if err != nil {
		return fmt.Fprintln(os.Stderr, message+":", err)
	}
	return
}

// PrintNonError writes info message to os.Stdout only when err is nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintNonError(err error, message string) (n int, errr error) {
	if err == nil {
		return fmt.Println(message)
	}
	return
}

// PrintErrorf formats according to a format specifier and writes error message to os.StdErr only when err is not nil.
// A newline is always appended. It returns the number of bytes written and any write error encountered.
// Please follow conventions for error message e.g. start with lowercase, don't end with period etc.
func (logger Logger) PrintErrorf(err error, message string, a ...interface{}) (n int, errr error) {
	if err != nil {
		a = append(a, err)
		return fmt.Fprintf(os.Stderr, message+":%+v\n", a...)
	}
	return
}

// PrintNonErrorf formats according to a format specifier and writes info message to os.Stdout only when err is nil.
// A newline is always appended. It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintNonErrorf(err error, message string, a ...interface{}) (n int, errr error) {
	if err == nil {
		return fmt.Printf(message+"\n", a...)
	}
	return
}

// WritefDebugMessageIfError formats according to a format specifier and writes to w only when DebugEnabled is true and err is not nil.
//  A newline is always appended. It returns the number of bytes written and any write error encountered.
func (logger Logger) WritefDebugMessageIfError(err error, w io.Writer, format string, a ...interface{}) (n int, errr error) {
	if err != nil {
		a = append(a, err)
		return logger.WritefDebugMessage(w, format+": %+v", a...)
	}
	return
}

// PrintfDebugMessageIfError formats according to a format specifier and writes to standard output only when DebugEnabled is true and err is not nil.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintfDebugMessageIfError(err error, format string, a ...interface{}) (n int, errr error) {
	return logger.WritefDebugMessageIfError(err, os.Stderr, format, a...)
}

// WritelnDebugMessageIfError formats using the default formats for its operands and writes to w only when DebugEnabled and err is not nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) WritelnDebugMessageIfError(err error, w io.Writer, a ...interface{}) (n int, errr error) {
	if err != nil {
		return logger.WritefDebugMessage(w, fmt.Sprintln(a...)+": %+v", err)
	}
	return
}

// PrintlnDebugMessageIfError formats using the default formats for its operands and writes to standard output only when DebugEnabled and err is not nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintlnDebugMessageIfError(err error, a ...interface{}) (n int, errr error) {
	return logger.WritelnDebugMessageIfError(err, os.Stderr, a...)
}

// WritefDebugMessageIfNotError formats according to a format specifier and writes to w only when DebugEnabled is true and err is nil.
//  A newline is always appended. It returns the number of bytes written and any write error encountered.
func (logger Logger) WritefDebugMessageIfNotError(err error, w io.Writer, format string, a ...interface{}) (n int, errr error) {
	if err == nil {
		return logger.WritefDebugMessage(w, format+": %+v", a...)
	}
	return
}

// PrintfDebugMessageIfNotError formats according to a format specifier and writes to standard output only when DebugEnabled is true and err is nil.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintfDebugMessageIfNotError(err error, format string, a ...interface{}) (n int, errr error) {
	return logger.WritefDebugMessageIfNotError(err, os.Stderr, format, a...)
}

// WritelnDebugMessageIfNotError formats using the default formats for its operands and writes to w only when DebugEnabled and err is nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) WritelnDebugMessageIfNotError(err error, w io.Writer, a ...interface{}) (n int, errr error) {
	if err == nil {
		return logger.WritelnDebugMessage(w, a...)
	}
	return
}

// PrintlnDebugMessageIfNotError formats using the default formats for its operands and writes to standard output only when DebugEnabled and err is nil.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (logger Logger) PrintlnDebugMessageIfNotError(err error, a ...interface{}) (n int, errr error) {
	return logger.WritelnDebugMessageIfNotError(err, os.Stderr, a...)
}
