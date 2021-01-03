package service

import (
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
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
	// 早上8:02
	_, err := s.cronMgr.AddFunc("2 8 * * *", s.MorningNoticeTask)
	if err != nil {
		log.Error().Msgf("add cron func err: %+v", err)
		return err
	}
	return nil
}

func (s *CronService) MorningNoticeTask() {
	log.Info().Msg("start morningNotice")
	loc, err := DefaultLocationService.GetMyLocation()
	if err != nil {
		log.Error().Msgf("GetMyLocation err: %+v", err)
		return
	}
	log.Info().Msgf("loc: %+v", loc)
	weather, err := DefaultWeatherService.GetWeatherByCity(loc.City)
	if err != nil {
		log.Error().Msgf("GetWeatherByCity err: %+v", err)
		return
	}
	log.Info().Msgf("weather: %+v", weather)
	buf := make([]byte, 0, 16)
	buf = append(buf, weather.City.Secondaryname...)
	buf = append(buf, ","...)
	buf = append(buf, weather.Condition.Condition...)
	buf = append(buf, ","...)
	buf = append(buf, weather.Condition.Tips...)
	buf = append(buf, weather.Sfc.Notice...)
	req := &api.TtsRequest{
		Text: string(buf),
		BaseRequest: api.BaseRequest{
			TraceID: uuid.New().String(),
		},
	}
	DefaultTtsService.TextToVoice(req)
}
