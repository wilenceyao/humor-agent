package service

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/config"
	"strconv"
)

var DefaultMqttService *MqttService

type MqttService struct {
}

func InitMqttService() error {
	DefaultMqttService = &MqttService{}
	return DefaultMqttService.startMqttClient()
}

func (s *MqttService) startMqttClient() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Config.Mqtt.Ip, config.Config.Mqtt.Port))
	opts.SetClientID(config.Config.ID)
	opts.SetUsername(config.Config.Mqtt.Username)
	opts.SetPassword(config.Config.Mqtt.Password)
	opts.SetDefaultPublishHandler(s.mqttPublishHandler)
	opts.OnConnect = s.onMqttConnect
	opts.OnConnectionLost = s.onMqttConnectionLost
	client := mqtt.NewClient(opts)
	return s.mqttConnect(client)
}

func (s *MqttService) onMqttConnect(c mqtt.Client) {
	topic := fmt.Sprintf("device/%s", config.Config.ID)
	token := c.Subscribe(topic, 1, nil)
	token.Wait()
	log.Info().Msgf("Subscribed to topic: %s", topic)
}

func (s *MqttService) onMqttConnectionLost(c mqtt.Client, err error) {
	log.Error().Msgf("Connect lost: %+v", err)
	s.mqttConnect(c)
}

func (s *MqttService) mqttPublishHandler(c mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())

	req := &api.TtsRequest{
		BaseRequest: api.BaseRequest{
			TraceID: strconv.FormatUint((uint64)(msg.MessageID()), 10),
		},
		Text: string(msg.Payload()),
	}
	DefaultTtsService.TextToVoice(req)
}

func (s *MqttService) mqttConnect(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
