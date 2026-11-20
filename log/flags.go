package log

import (
	"flag"
	"github.com/sirupsen/logrus"
)

const (
	logLevel = "log.level"
)

type LoggingFlags struct {
	Level logrus.Level
}

var loggingFlags = &LoggingFlags{
	Level: logrus.InfoLevel,
}

func AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(logLevel, "info", "Minimal allowed log Level.")
}

//func Init(cfg *config.Config) *LoggingFlags {
//	levelstring := cfg.GetString(logLevel)
//	level, err := logrus.ParseLevel(levelstring)
//	if err != nil {
//		fmt.Printf("invalided log level %s, will be seted as info log level", levelstring)
//		return loggingFlags
//	}
//
//	loggingFlags.Level = level
//
//	return loggingFlags
//}
