package main

import (
	"fmt"
        "log/syslog"
        logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"

	"runtime"

	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var log *logrus.Logger

func PanicRecover() {
    if r := recover(); r != nil {
        logrus.Error("Internal error: %v", r)
        buf := make([]byte, 1<<16)
        stackSize := runtime.Stack(buf, true)
        log.Error("--------------------------------------------------------------------------------")
        log.Error(fmt.Sprintf("Internal error: %s\n", string(buf[0:stackSize])))
        log.Error("--------------------------------------------------------------------------------")
        fmt.Printf("--------------------------------------------------------------------------------\n")
        fmt.Printf(fmt.Sprintf("Internal error: %s\n", string(buf[0:stackSize])))
        fmt.Printf("--------------------------------------------------------------------------------\n")

        }
}

func InitLogrus() *logrus.Logger {

	log := logrus.New()
        log.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
	if err == nil {
		log.Hooks.Add(hook)
	}

	log.AddHook(&writer.Hook{ // Send logs with level higher than warning to stderr
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{ // Send info and debug logs to stdout
		Writer: os.Stdout,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.TraceLevel,
			logrus.DebugLevel,
		},
	})

	log.SetReportCaller(true)

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetLevel(logrus.TraceLevel)

	return log
}

func divideByZero() {
    defer PanicRecover()
    fmt.Println(divide(1, 0))
}

func divide(a, b int) int {
    if b == 0 {
        panic(fmt.Sprintln("DEV BY ZERO",a,b))
    }
    return a / b
}


func main() {

	log = InitLogrus()

        defer PanicRecover()


	log.Info("This will go to stdout")
	log.Warn("This will go to stderr")


	fmt.Printf("%#v\n",log)

	log.Trace("Trace message")
	log.Info("Info message")
	log.Warn("Warn message")
	log.Error("Error message")
//	log.Fatal("Fatal message")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Info("A walrus appears????????")
	log.Info("A walrus appears!!!!!!")

        log.Info("Info message")
        log.Info("Info message")
        log.Info("Info message")
        log.Info("Info message")

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	contextLogger := log.WithFields(logrus.Fields{
		"common": "this is a common field",
		"other": "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")

        divideByZero()

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")

}
