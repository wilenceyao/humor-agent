package internal

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/config"
)

func (a *HumorAgent) startMqttClient() error {
	// netInterface, err := net.InterfaceByName(currentNetworkHardwareName)
	// macAddress := netInterface.HardwareAddr
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Config.Mqtt.Ip, config.Config.Mqtt.Port))
	opts.SetClientID(id)
	opts.SetUsername(config.Config.Mqtt.Username)
	opts.SetPassword(config.Config.Mqtt.Password)
	opts.SetDefaultPublishHandler(a.mqttPublishHandler)
	opts.OnConnect = a.onMqttConnect
	opts.OnConnectionLost = a.onMqttConnectionLost
	client := mqtt.NewClient(opts)
	return a.mqttConnect(client)
}

func (a *HumorAgent) onMqttConnect(c mqtt.Client) {
	topic := id
	token := c.Subscribe(id, 1, nil)
	token.Wait()
	log.Info().Msgf("Subscribed to topic: %s", topic)
}

func (a *HumorAgent) onMqttConnectionLost(c mqtt.Client, err error) {
	log.Error().Msgf("Connect lost: %+v", err)
	a.mqttConnect(c)
}

func (a *HumorAgent) mqttPublishHandler(c mqtt.Client, msg mqtt.Message) {
	log.Info().Msgf("Received message: %s from topic: %s", msg.Payload(), msg.Topic())
}

func (a *HumorAgent) mqttConnect(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
