package service

import (
	"context"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
)

type AgentServiceImpl struct {
}

func (s *AgentServiceImpl) Tts(ctx context.Context, req *agentapi.TtsRequest) (*agentapi.TtsResponse, error) {
	res := &agentapi.TtsResponse{
		Response: &common.BaseResponse{},
	}
	res.Response = &common.BaseResponse{}
	DefaultTtsService.TextToVoice(req, res)
	return res, nil
}

func (s *AgentServiceImpl) Weather(ctx context.Context, req *agentapi.WeatherRequest) (*agentapi.WeatherResponse, error) {
	res := &agentapi.WeatherResponse{
		Response: &common.BaseResponse{},
	}
	DefaultWeatherService.LocalWeather(req, res)
	return res, nil
}
