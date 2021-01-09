package service

import (
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
)

var DefaultCronService *CronService

type CronService struct {
	cronMgr *cron.Cron
}

func InitCronService() error {
	DefaultCronService = &CronService{
		cronMgr: cron.New(),
	}
	return DefaultCronService.addCron()
}

func (s *CronService) addCron() error {
	// 早上07:35
	_, err := s.cronMgr.AddFunc("35 7 * * *", s.localWeather)
	if err != nil {
		log.Error().Msgf("add cron func err: %+v", err)
		return err
	}
	s.cronMgr.Start()
	return nil
}

func (s *CronService) localWeather() {
	req := &agentapi.WeatherRequest{
		Request: &common.BaseRequest{
			RequestID: uuid.New().String(),
		},
	}
	res := &agentapi.WeatherResponse{
		Response: &common.BaseResponse{},
	}
	DefaultWeatherService.LocalWeather(req, res)
}
