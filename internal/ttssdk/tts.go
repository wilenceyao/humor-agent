package ttssdk

import (
	"bytes"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/config"
	"github.com/wilenceyao/humor-agent/internal/player"
)

// 腾讯云API文档: https://cloud.tencent.com/document/product/1073/37995
var client *tts.Client
var (
	modelType int64  = 1
	codec     string = "mp3"
)

func Init() {

	credential := common.NewCredential(
		config.Config.QCloud.SecretId,
		config.Config.QCloud.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tts.tencentcloudapi.com"
	client, _ = tts.NewClient(credential, "ap-guangzhou", cpf)
}

func TextToVoice(req *api.TtsRequest) *api.TtsReply {
	reply := &api.TtsReply{
		Reply: &api.BaseReply{},
	}
	request := tts.NewTextToVoiceRequest()
	request.Text = common.StringPtr(req.Text)
	request.VoiceType = &config.Config.QCloud.TtsConfig.VoiceType
	request.ModelType = &modelType
	request.SessionId = &req.Req.TraceId
	request.Codec = &codec
	response, err := client.TextToVoice(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			log.Error().Msgf("TencentCloudSDKError : %+v", err)
		}
		reply.Reply.Code = api.ErrorCode_EXTERNAL_ERROR
		return reply
	}
	arr, err := base64.StdEncoding.DecodeString(*response.Response.Audio)
	if err != nil {
		log.Error().Msgf("base64 decode err: %+v", err)
		reply.Reply.Code = api.ErrorCode_INTERNAL_ERROR
		return reply
	}
	id := uuid.New().String()
	audio := player.Audio{
		R:     bytes.NewReader(arr),
		Title: id,
	}
	player.Player.Enqueue(audio)
	reply.Id = id
	return reply
}
