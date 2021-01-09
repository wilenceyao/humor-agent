package service

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	qcloudcommon "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
	"github.com/wilenceyao/humor-agent/config"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
)

var DefaultTtsService *TtsService

type TtsService struct {
	client *tts.Client
}

// 腾讯云API文档: https://cloud.tencent.com/document/product/1073/37995

var (
	modelType int64  = 1
	codec     string = "mp3"
)

func InitTtsService() {
	credential := qcloudcommon.NewCredential(
		config.Config.QCloud.SecretId,
		config.Config.QCloud.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tts.tencentcloudapi.com"
	qcloudTtsClient, _ := tts.NewClient(credential, "ap-guangzhou", cpf)
	DefaultTtsService = &TtsService{
		client: qcloudTtsClient,
	}
}

func (s *TtsService) TextToVoice(req *agentapi.TtsRequest, res *agentapi.TtsResponse) {
	request := tts.NewTextToVoiceRequest()
	request.Text = qcloudcommon.StringPtr(req.Text)
	request.VoiceType = &config.Config.QCloud.TtsConfig.VoiceType
	request.ModelType = &modelType
	request.SessionId = &req.Request.RequestID
	request.Codec = &codec
	response, err := s.client.TextToVoice(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			log.Error().Msgf("TencentCloudSDKError : %+v", err)
		}
		res.Response.Code = common.ErrorCode_EXTERNAL_ERROR
		return
	}
	arr, err := base64.StdEncoding.DecodeString(*response.Response.Audio)
	if err != nil {
		log.Error().Msgf("base64 decode err: %+v", err)
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
		return
	}
	id := uuid.New().String()
	audio := Audio{
		R:     bytes.NewReader(arr),
		Title: id,
	}
	DefaultTtsPlayer.player.Enqueue(audio)
}
