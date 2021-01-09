package service

import (
	"context"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/config"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
	"github.com/wilenceyao/humors"
)

var DefaultMqttService *MqttService

type MqttService struct {
	humorSys *humors.Humors
}

func InitMqttService() error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", config.Config.Mqtt.IP, config.Config.Mqtt.Port))
	opts.SetUsername(config.Config.Mqtt.Username)
	opts.SetPassword(config.Config.Mqtt.Password)
	opts.SetClientID(config.Config.Mqtt.ClientID)
	h, err := humors.NewHumors(opts)
	if err != nil {
		log.Error().Msgf("init humors err: %v", err)
		return err
	}
	h.InitServant("")
	DefaultMqttService = &MqttService{
		humorSys: h,
	}
	h.Servant.RegisterFun(int32(agentapi.Action_TTS), agentapi.TtsRequest{}, agentapi.TtsResponse{},
		DefaultMqttService.tts)

	return nil
}

func (s *MqttService) tts(ctx context.Context, reqObj interface{}, resObj interface{}) {
	req := reqObj.(*agentapi.TtsRequest)
	res := resObj.(*agentapi.TtsResponse)
	res.Response = &common.BaseResponse{}
	DefaultTtsService.TextToVoice(req, res)
}
