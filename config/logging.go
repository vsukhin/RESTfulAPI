package config

import (
	logging "github.com/op/go-logging"
	/*
		gelf "github.com/probkiizokna/go-gelf"
		logging_gelf "github.com/probkiizokna/go-logging-gelf"
	*/
	"errors"
	"os"
	"path"
	"syscall"
)

type Logger interface {
	Info(query string, args ...interface{})
	Warning(query string, args ...interface{})
	Error(query string, args ...interface{})
	Fatalf(query string, args ...interface{})
}

const (
	LOGGING_CONSOLE = "console"
	LOGGING_SYSLOG  = "syslog"
	LOGGING_FILE    = "file"
	LOGGING_GRAYLOG = "graylog"
)

var (
	Filebackend *FileBackend
	logger      Logger = logging.MustGetLogger("config")
)

func InitLogging() (err error) {
	var backends []logging.Backend
	for _, mode := range Configuration.Logger.Mode {
		level := logging.DEBUG
		modelevel := string(Configuration.Logger.Levels[ModeName(mode)])
		if modelevel != "" {
			level, err = logging.LogLevel(modelevel)
			if err != nil {
				logger.Fatalf("Can't recognize mode level %v", err)
				return err
			}
		}
		switch mode {
		case LOGGING_CONSOLE:
			stdoutbackend := logging.NewLogBackend(os.Stdout, "", 0)
			stdoutbackend.Color = true
			leveledstdout := logging.AddModuleLevel(stdoutbackend)
			leveledstdout.SetLevel(level, "")
			backends = append(backends, stdoutbackend)
		case LOGGING_SYSLOG:
			syslogbackend, err := logging.NewSyslogBackend(CONFIG_APPLICATION)
			if err != nil {
				logger.Fatalf("Can't initiate syslog backend %v", err)
				return err
			}
			leveledsyslog := logging.AddModuleLevel(syslogbackend)
			leveledsyslog.SetLevel(level, "")
			backends = append(backends, leveledsyslog)
		case LOGGING_FILE:
			file, err := os.OpenFile(Configuration.Logger.File, syscall.O_APPEND|syscall.O_CREAT|syscall.O_WRONLY, 0666)
			if err != nil {
				logger.Fatalf("Can't initiate filelog backend%v", err)
				return err
			}
			Filebackend = NewFileBackend(file)
			leveledfile := logging.AddModuleLevel(Filebackend)
			leveledfile.SetLevel(level, "")
			backends = append(backends, leveledfile)
		case LOGGING_GRAYLOG:
			/*
				gelfClient := gelf.MustUdpClient(
					Configuration.Logger.Graylog.Host,
					Configuration.Logger.Graylog.Port,
					Configuration.Logger.Graylog.ChunkSize,
					Configuration.Logger.Graylog.Compression,
				)
				gelfBacked := logging_gelf.NewGelfBackend(gelfClient, mustHostname(), mustApplicationName(true))
				leveledgelf = logging.AddModuleLevel(gelfbackend)
				leveledgelf.SetLevel(level, "")
				backends = append(backends, leveledgelf)
			*/
		default:
			logger.Fatalf("Uknown logging mode %v", mode)
			return errors.New("Uknown logging mode")
		}
		logging.SetBackend(backends...)
		logFormatter := logging.MustStringFormatter(Configuration.Logger.Format)
		logging.SetFormatter(logFormatter)
	}

	return nil
}

func ShutdownLogging() {
	if Filebackend.File != nil {
		err := Filebackend.File.Close()
		if err != nil {
			logger.Error("Log file close: %v", err)
		}
	}
}

func mustHostname() string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	} else {
		logger.Fatalf("Can't recognize host name %v", err)
		return ""
	}
}

func mustApplicationName(short bool) string {
	if applicationName := os.Args[0]; applicationName != "" {
		if short {
			return path.Base(applicationName)
		} else {
			return applicationName
		}
	} else {
		logger.Fatalf("Can't detect application name")
		return ""
	}
}
