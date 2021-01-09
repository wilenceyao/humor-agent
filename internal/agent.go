package internal

import (
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/config"
	"github.com/wilenceyao/humor-agent/internal/admin"
	"github.com/wilenceyao/humor-agent/internal/service"
	"gopkg.in/natefinch/lumberjack.v2"
	stdlog "log"
)

type HumorAgent struct {
}

func (a *HumorAgent) Start() error {
	var err error
	err = config.Init("config.json")
	if err != nil {
		log.Error().Msgf("config init err: %+v", err)
		return err
	}
	a.setupLog()
	log.Info().Msg("Starting")
	err = service.InitTtsPlayer()
	if err != nil {
		log.Error().Msgf("init tts player err: %+v", err)
		return err
	}
	service.InitTtsService()
	err = service.InitMqttService()
	if err != nil {
		log.Error().Msgf("init mqtt client err: %+v", err)
		return err
	}
	err = service.InitCronService()
	if err != nil {
		log.Error().Msgf("init cron service err: %+v", err)
		return err
	}
	return admin.RunAdminServer()
}

func (a *HumorAgent) setupLog() {
	h := &lumberjack.Logger{
		Filename:   config.Config.LogFile,
		MaxSize:    100,  // megabytes
		MaxBackups: 10,   // 最多50个日志文件，因而只保留49个旧日志备份
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	}
	stdlog.SetOutput(h)
	log.Logger = log.With().Caller().Logger().Output(h)
}
