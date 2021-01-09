package service

import (
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/rs/zerolog/log"
	"io"
)

var DefaultTtsPlayer *TtsPlayer

type TtsPlayer struct {
	player *AudioPlayer
}

type Audio struct {
	R     io.Reader
	Title string
}

type AudioPlayer struct {
	c           chan Audio
	innerPlayer *oto.Player
	innerCtx    *oto.Context
}

func InitTtsPlayer() error {
	ctx, err := oto.NewContext(16000, 2, 2, 8192)
	if err != nil {
		return err
	}
	player := &AudioPlayer{
		c:           make(chan Audio, 10),
		innerCtx:    ctx,
		innerPlayer: ctx.NewPlayer(),
	}
	DefaultTtsPlayer = &TtsPlayer{
		player: player,
	}
	go DefaultTtsPlayer.player.playCoroutine()
	return nil
}

func (p *AudioPlayer) Close() {
	p.innerPlayer.Close()
	p.innerCtx.Close()
}

func (p *AudioPlayer) playCoroutine() {
	var err error
	for {
		a := <-p.c
		log.Info().Msgf("start play %s", a.Title)
		err = p.play(a)
		if err != nil {
			log.Error().Msgf("play audio %s err: %+v", a.Title, err)
		}
		log.Info().Msgf("end play %s", a.Title)
	}
}

func (p *AudioPlayer) Enqueue(a Audio) {
	p.c <- a
}

func (p *AudioPlayer) play(a Audio) error {
	log.Info().Msgf("playing %s", a.Title)
	decoder, err := mp3.NewDecoder(a.R)
	if err != nil {
		return err
	}
	log.Info().Msgf("file length: %d", decoder.Length())
	if _, err := io.Copy(p.innerPlayer, decoder); err != nil {
		return err
	}
	return nil
}
