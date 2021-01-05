package service

import (
	"context"
	"github.com/eclipse/paho.golang/paho"
	"github.com/eclipse/paho.golang/paho/extensions/rpc"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/config"
	mqttapi "github.com/wilenceyao/humor-api/api/mqtt"
	humorutil "github.com/wilenceyao/humor-api/pkg/util"
	"google.golang.org/protobuf/proto"
)

var DefaultMqttService *MqttService

type MqttService struct {
	MqttRpcReqHandler *rpc.Handler
	MqttClient        *paho.Client
}

func InitMqttService() error {
	DefaultMqttService = &MqttService{}
	return DefaultMqttService.startMqttClient()
}

func (s *MqttService) rpcDispatcher(m *paho.Publish) {
	log.Info().Msgf("Received message: %s from topic: %s", m.Payload, m.Topic)
	mqttMsg := &mqttapi.Message{}
	err := proto.Unmarshal(m.Payload, mqttMsg)
	if err != nil {
		log.Error().Msgf("mqtt msg unmarshal err: %+v", err)
		return
	}
	switch mqttMsg.Action {
	case mqttapi.Action_TTS:
		ttsReq := &mqttapi.TtsRequest{}
		proto.Unmarshal(mqttMsg.Payload, ttsReq)
		req := &api.TtsRequest{
			BaseRequest: api.BaseRequest{
				TraceID: mqttMsg.TraceID,
			},
			Text: ttsReq.Text,
		}
		DefaultTtsService.TextToVoice(req)

		ttsRes := &mqttapi.TtsReply{
			Reply: &mqttapi.BaseReply{
				Code: mqttapi.ErrorCode_SUCCESS,
			},
		}
		ttsResBtArr, _ := proto.Marshal(ttsRes)
		mqttResMsg := &mqttapi.Message{
			Action:  mqttMsg.Action,
			TraceID: mqttMsg.TraceID,
			Payload: ttsResBtArr,
		}
		mqttResBtArr, _ := proto.Marshal(mqttResMsg)
		_, err := s.MqttClient.Publish(context.Background(), &paho.Publish{
			Properties: &paho.PublishProperties{
				CorrelationData: m.Properties.CorrelationData,
			},
			Topic:   m.Properties.ResponseTopic,
			Payload: mqttResBtArr,
		})
		if err != nil {
			log.Error().Msgf("response %s err: %+v", mqttMsg.TraceID, err)
		}

	}
}

func (s *MqttService) startMqttClient() error {
	rpcConfig := &humorutil.MqttRpcConfig{
		ClientID:    config.Config.ID,
		IP:          config.Config.Mqtt.IP,
		Port:        config.Config.Mqtt.Port,
		Username:    config.Config.Mqtt.Username,
		Password:    config.Config.Mqtt.Password,
		RecvHandler: s.rpcDispatcher,
	}
	var err error
	s.MqttClient, s.MqttRpcReqHandler, err = humorutil.NewMqttRpcHandler(rpcConfig)
	if err != nil {
		log.Error().Msgf("InitMqttRpc err: %+v", err)
		return err
	}
	return nil
}

func (s *MqttService) mqttPublishHandler(c mqtt.Client, msg mqtt.Message) {

}

func (s *MqttService) mqttConnect(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
