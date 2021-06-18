package minio

import (
	"ceres/pkg/config"
	"context"
	"time"

	"github.com/gotomicro/ego/core/elog"
	"github.com/minio/minio-go/v7"
)

func PreSignUpload(client *minio.Client, file string) (url string, err error) {

	u, err := client.PresignedPutObject(
		context.TODO(),
		config.Minio.Bucket,
		file,
		time.Minute*10,
	)
	if err != nil {
		elog.Errorf("PreSign the upload request failed %v", err)
		return
	}
	url = u.String()
	return
}
