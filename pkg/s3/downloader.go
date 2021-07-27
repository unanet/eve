package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/unanet/go/pkg/errors"
)

type Downloader struct {
	s3 *s3manager.Downloader
}

func NewDownloader(sess *session.Session) *Downloader {
	return &Downloader{
		s3: s3manager.NewDownloader(sess),
	}
}

func (u Downloader) Download(ctx context.Context, location *Location) ([]byte, error) {
	buf := &aws.WriteAtBuffer{}

	_, err := u.s3.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(location.Bucket),
		Key:    aws.String(location.Key),
	})
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return buf.Bytes(), nil
}
