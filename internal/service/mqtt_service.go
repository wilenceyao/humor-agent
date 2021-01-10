package service

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/config"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humors"
)

var DefaultMqttService *MqttService

type MqttService struct {
	humorSys *humors.Humors
}

func InitMqttService() error {
	mqttOpts := MQTT.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("%s:%d", config.Config.Mqtt.IP, config.Config.Mqtt.Port))
	mqttOpts.SetUsername(config.Config.Mqtt.Username)
	mqttOpts.SetPassword(config.Config.Mqtt.Password)
	mqttOpts.SetClientID(config.Config.Mqtt.ClientID)
	rpcOpts := &humors.RPCOptions{
		Timeout: 2000,
	}
	opts := humors.Options{
		MQTTOpts: mqttOpts,
		RPCOpts:  rpcOpts,
	}
	h, err := humors.NewHumors(opts)
	if err != nil {
		log.Error().Msgf("init humors err: %v", err)
		return err
	}
	agentapi.RegisterAgentServiceServer(h, &AgentServiceImpl{})
	DefaultMqttService = &MqttService{
		humorSys: h,
	}
	return nil
}
