package logs

import (
	"os"
	config "starland-account/configs"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogging(cfg *config.Config) {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logfile/log.log",  
			MaxSize:    10,                 
			MaxBackups: 5,                 
			MaxAge:     30,                
			Compress:   true,               
		})),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
	zap.S().Infof("InitLogging: zap initialized")
}
