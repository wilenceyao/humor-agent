package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/config"
	"github.com/wilenceyao/humor-agent/internal/cmasdk"
	"github.com/wilenceyao/humor-agent/internal/location"
	"github.com/wilenceyao/humor-agent/internal/player"
	"github.com/wilenceyao/humor-agent/internal/ttssdk"
	"google.golang.org/grpc"
	"net"
)

type HumorAgent struct {
	api.UnimplementedHumorAgentServer
	cronMgr *cron.Cron
}

var port int = 8080
var id string = "82:a9:10:86:38:01"

func (a *HumorAgent) Start() error {
	var err error
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("Starting")
	err = config.Init("config.json")
	if err != nil {
		log.Error().Msgf("config init err: %+v", err)
		return err
	}
	ttssdk.Init()
	if err = player.Init(); err != nil {
		log.Error().Msgf("player init err: %+v", err)
		return err
	}
	a.cronMgr = cron.New()
	a.addCron()
	err = a.startMqttClient()
	if err != nil {
		log.Error().Msgf("start mqtt client err: %+v", err)
		return err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error().Msgf("failed to listen: %v", err)
		return err
	}
	log.Info().Msg("Started")
	s := grpc.NewServer()
	api.RegisterHumorAgentServer(s, a)
	if err := s.Serve(lis); err != nil {
		log.Error().Msgf("failed to serve: %v", err)
		return err
	}
	return nil
}

func (a *HumorAgent) Tts(ctx context.Context, req *api.TtsRequest) (*api.TtsReply, error) {
	reply := ttssdk.TextToVoice(req)
	log.Info().Msgf("[API] Tts, req: %+v, res: %+v", req, reply)
	return reply, nil
}

func (a *HumorAgent) addCron() error {
	// 早上7:55
	_, err := a.cronMgr.AddFunc("55 7 * * *", a.morningNotice)
	if err != nil {
		log.Error().Msgf("add cron func err: %+v", err)
		return err
	}
	return nil
}

func (a *HumorAgent) morningNotice() {
	log.Info().Msg("start morningNotice")
	loc, err := location.GetMyLocation()
	if err != nil {
		log.Error().Msgf("GetMyLocation err: %+v", err)
		return
	}
	log.Info().Msgf("loc: %+v", loc)
	weather, err := cmasdk.GetWeatherByCity(loc.City)
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
	req := &api.TtsRequest{Text: string(buf), Req: &api.BaseRequest{
		TraceId: uuid.New().String(),
	}}
	ttssdk.TextToVoice(req)
}
