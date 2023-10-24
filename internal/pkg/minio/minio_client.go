package minio

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mime/multipart"
	"path/filepath"
)

type Minio struct {
	Client     *minio.Client
	logger     *logrus.Logger
	BucketName string
	Host       string
}

type Client interface {
	SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, objectName string) error
}

type MinioConfig struct {
	Host            string
	BucketName      string
	AccessKeyID     string
	SecretAccessKey string
	Location        string
}

func InitConfig(vp *viper.Viper) MinioConfig {
	config := MinioConfig{
		Host:            vp.GetString("minio.host"),
		BucketName:      vp.GetString("minio.bucketName"),
		AccessKeyID:     vp.GetString("minio.accessKeyId"),
		SecretAccessKey: vp.GetString("minio.SecretAccessKey"),
		Location:        vp.GetString("location"),
	}

	return config
}

func NewMinioClient(ctx context.Context, config MinioConfig, logger *logrus.Logger) (Client, error) {
	minioClient, err := minio.New(config.Host, &minio.Options{
		Creds: credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
	})

	if err != nil {
		logger.Fatalf("error: %s", err)
		return nil, err
	}

	if err = minioClient.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{Region: config.Location}); err != nil {
		fmt.Println(err)
		if exists, err := minioClient.BucketExists(ctx, config.BucketName); err == nil && exists {
			logger.Println("We already own ", config.BucketName)
		} else if err != nil {
			logger.Fatalln(err)
		}
	}

	logger.Println("Successfully created ", config.BucketName)

	return &Minio{Client: minioClient, BucketName: config.BucketName, logger: logger, Host: config.Host}, nil
}

func (m *Minio) SaveImage(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	objectName := uuid.New().String() + filepath.Ext(header.Filename)

	if _, err := m.Client.PutObject(ctx, m.BucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	}); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s/%s", m.Host, m.BucketName, objectName), nil
}

func (m *Minio) DeleteImage(ctx context.Context, objectName string) error {
	return m.Client.RemoveObject(ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
}
