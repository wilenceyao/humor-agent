package service

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	agentapi "github.com/wilenceyao/humor-api/agent/humor"
	"github.com/wilenceyao/humor-api/common"
	"gocv.io/x/gocv"
)

type PhotoService struct {
	deviceID int
	fileExt  gocv.FileExt
}

var DefaultPhotoService *PhotoService

func InitPhotoService() {
	DefaultPhotoService = &PhotoService{
		deviceID: 0,
		fileExt:  gocv.JPEGFileExt,
	}
}

func (s *PhotoService) TakeAndUpload(ctx context.Context, req *agentapi.TakePhotoRequest,
	res *agentapi.TakePhotoResponse) {
	mat := gocv.NewMat()
	defer mat.Close()
	err := s.takePhoto(&mat)
	if err != nil {
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
		res.Response.Msg = err.Error()
		return
	}
	btArr, err := gocv.IMEncode(s.fileExt, mat)
	if err != nil {
		res.Response.Code = common.ErrorCode_INTERNAL_ERROR
		res.Response.Msg = err.Error()
		return
	}
	photoID := req.Id
	name := fmt.Sprintf("%s%s", photoID, s.fileExt)
	// 拍照比较耗时，导致ctx可能已经超时
	err = DefaultCOSService.PutPhoto(ctx, name, btArr)
	if err != nil {
		res.Response.Code = common.ErrorCode_EXTERNAL_ERROR
		res.Response.Msg = err.Error()
	}
}

func (s *PhotoService) takePhoto(mat *gocv.Mat) error {
	webcam, err := gocv.OpenVideoCapture(s.deviceID)
	if err != nil {
		log.Error().Msgf("open video capture device err: %v", err)
		return err
	}
	defer webcam.Close()
	if ok := webcam.Read(mat); !ok {
		log.Error().Msgf("read capture device err: %v", err)
		return err
	}
	if mat.Empty() {
		err = fmt.Errorf("no image on capture device")
		log.Err(err)
		return err
	}
	return nil
}
