package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"gitlab.unanet.io/devops/eve/pkg/errors"
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

func (u Uploader) UploadText(ctx context.Context, key string, body string) (*Location, error) {
	bodyReader := strings.NewReader(body)
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
