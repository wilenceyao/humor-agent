package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/wilenceyao/humor-agent/api"
	"github.com/wilenceyao/humor-agent/internal/service"
	"net/http"
)

type Service struct {
}

func (a *Service) Tts(c *gin.Context) {
	var req api.TtsRequest
	res := &api.TtsResponse{}
	if err := c.ShouldBindJSON(&req); err != nil {
		res.Code = api.INVALID_PARAMETERS
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res = service.DefaultTtsService.TextToVoice(&req)
	c.JSON(http.StatusOK, res)
}

func (a *Service) MorningNotice(c *gin.Context) {
	res := &api.MorningNoticeResponse{}
	service.DefaultCronService.MorningNoticeTask()
	c.JSON(http.StatusOK, res)
}
