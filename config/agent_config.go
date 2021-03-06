package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
)

var Config *AgentConfig

type AgentConfig struct {
	QCloud QCloudConfig
	Mqtt   MqttConfig
	// 本地管控服务
	Admin   AdminConfig
	LogFile string
}

type AdminConfig struct {
	Port int
}

type QCloudConfig struct {
	SecretId  string
	SecretKey string
	TtsConfig QCloudTtsConfig
	CosConfig QCloudCosConfig
}

type MqttConfig struct {
	IP       string
	Port     uint
	Username string
	Password string
	ClientID string
}

type QCloudTtsConfig struct {
	VoiceType int64
}

type QCloudCosConfig struct {
	Bucket string
	Region string
	PhotoDir string
}

func Init(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Error().Msgf("read config file err: %+v", err)
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = &AgentConfig{}
	err = decoder.Decode(Config)
	if err != nil {
		log.Error().Msgf("decode config file err: %+v", err)
		return err
	}
	return nil
}
