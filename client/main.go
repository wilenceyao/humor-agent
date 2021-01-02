package main

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wilenceyao/humor-agent/api"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Msgf("did not connect: %v", err)
		return
	}
	defer conn.Close()
	c := api.NewHumorAgentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Tts(ctx, &api.TtsRequest{
		Req: &api.BaseRequest{
			TraceId: uuid.New().String(),
		},
		Text: "我都不好意思听",
	})
	if err != nil {
		log.Error().Msgf("call rpc err: %+v", err)
		return
	}
	if r.Reply.Code != api.ErrorCode_SUCCESS {
		log.Error().Msgf("rpc failed: %+v", r.Reply)
		return
	}
	log.Info().Msgf("msgId : %s", r.Id)
}
