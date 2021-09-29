package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"}, //loggin to stdout , In this way we are going to rely on logstash to push all this log to elastic search
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	//Now we need to define an error because when we try to  init the logging system we get as a result, a pointer to zaplogger(*zap.Logger) and an error
	//we are initializing the Log *zap.Logger
	var err error
	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}

}

func Info(msg string, tags ...zap.Field) {
	//log.Info(msg)
	log.Info(msg, tags...)
	log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.Error(msg, tags...)
	log.Sync()
}
