package internal

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/config"
	"github.com/wilenceyao/humor-agent/internal/player"
	"github.com/wilenceyao/humor-agent/internal/ttssdk"
	"google.golang.org/grpc"
	"net"
)

type HumorAgent struct {
	api.UnimplementedHumorAgentServer
}

var port int = 8080

func (a *HumorAgent) Start() error {
	var err error
	log.Info().Msg("Starting")
	config.Init("config.json")
	ttssdk.Init()
	if err = player.Init(); err != nil {
		log.Error().Msgf("player init err: %+v", err)
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
