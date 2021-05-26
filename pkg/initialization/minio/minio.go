package minio

import (
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/elog"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client
var Bucket string

/// minio configuration from the configuration files
type config struct {
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
}

func Init() (err error) {
	// init the minio client

	C := &config{}
	err = econf.UnmarshalKey("ceres.minio", C) //TODO：should change the logic to compile variables
	if err != nil {
		elog.Panicf("Parsing the minio configurations faild %v", err)
	}

	if C.Bucket == "" {
		elog.Panicf("Could not empty bucket %v", err)
	}
	Bucket = C.Bucket

	Client, err = minio.New(
		C.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(C.AccessKey, C.SecretKey, ""),
			Secure: true,
		},
	)

	if err != nil {
		elog.Panicf("Configure the client failed %v", err)
	}

	return
}
