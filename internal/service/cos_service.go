package service

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/wilenceyao/humor-agent/config"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
)

type COSService struct {
	cosClient *cos.Client
	photoDir  string
}

var DefaultCOSService *COSService

func InitCOSService() error {
	u, err := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		config.Config.QCloud.CosConfig.Bucket,
		config.Config.QCloud.CosConfig.Region))
	if err != nil {
		return err
	}
	b := &cos.BaseURL{BucketURL: u}
	// 1.永久密钥
	cosClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.Config.QCloud.SecretId,
			SecretKey: config.Config.QCloud.SecretKey,
		},
	})
	DefaultCOSService = &COSService{
		cosClient: cosClient,
		photoDir:  config.Config.QCloud.CosConfig.PhotoDir,
	}
	return nil
}

func (s *COSService) PutPhoto(ctx context.Context, name string, btArr []byte) error {
	name = fmt.Sprintf("%s/%s", s.photoDir, name)
	_, err := s.cosClient.Object.Put(ctx, name, bytes.NewReader(btArr), nil)
	if err != nil {
		log.Error().Msgf("put cos object %s err: %v", name, err)
	}
	return err
}
