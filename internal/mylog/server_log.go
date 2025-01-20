package mylog

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var (
	LOG_DEBUG = true
)
var CommonLogger zerolog.Logger
var FileLogger zerolog.Logger

type LogConfig struct {
	Level           zerolog.Level
	ConsoleEnabled  bool
	SaveFileEnabled bool
	FileName        string
}

func CreateLogger(config LogConfig) zerolog.Logger {
	writers := make([]io.Writer, 0)
	if config.ConsoleEnabled {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timeFormat}
		consoleWriter.FieldsExclude = []string{"caller"}
		//consoleWriter.FormatLevel = func(i interface{}) string {
		//    return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		//}
		//consoleWriter.FormatMessage = func(i interface{}) string {
		//    return fmt.Sprintf("%s", i)
		//}
		//consoleWriter.FormatFieldName = func(i interface{}) string {
		//    return fmt.Sprintf("%s:", i)
		//}
		//consoleWriter.FormatFieldValue = func(i interface{}) string {
		//    return fmt.Sprintf("%s;", i)
		//}

		writers = append(writers, consoleWriter)
	}
	if config.SaveFileEnabled {
		rollFile := &lumberjack.Logger{
			Filename: config.FileName,
			MaxSize:  100, // megabytes
			MaxAge:   1,   //days
			//Compress:   true, // disabled by default
		}
		writers = append(writers, rollFile)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	logger := zerolog.New(multi).Level(config.Level).With().Timestamp().Caller().Logger()
	return logger
}

func ManualInit() {
	zerolog.TimeFieldFormat = timeFormat
	level := zerolog.InfoLevel
	if LOG_DEBUG {
		level = zerolog.DebugLevel
	}
	CommonLogger = CreateLogger(LogConfig{
		Level:           level,
		ConsoleEnabled:  true,
		SaveFileEnabled: true,
		FileName:        "log/access.log",
	})

	FileLogger = CreateLogger(LogConfig{
		Level:           zerolog.DebugLevel,
		SaveFileEnabled: true,
		FileName:        "log/debug.log",
	})
}
func init() {
	ManualInit()
}
