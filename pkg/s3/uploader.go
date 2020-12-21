package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gitlab.unanet.io/devops/go/pkg/errors"
)

type Config struct {
	Bucket string
}

type Location struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Url    string `json:"url"`
}

type Uploader struct {
	Bucket string
	s3     *s3manager.Uploader
}

func NewUploader(sess *session.Session, config Config) *Uploader {
	return &Uploader{
		s3:     s3manager.NewUploader(sess),
		Bucket: config.Bucket,
	}
}

func (u Uploader) Upload(ctx context.Context, key string, body []byte) (*Location, error) {
	bodyReader := bytes.NewReader(body)
	result, err := u.s3.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   bodyReader,
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return &Location{
		Bucket: u.Bucket,
		Key:    key,
		Url:    result.Location,
	}, nil
}
