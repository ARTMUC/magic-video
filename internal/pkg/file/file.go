package file

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Reader interface {
	Read(path string) ([]byte, error)
}

type Writer interface {
	Write(path string, data []byte) error
}

type ReaderWriter interface {
	Reader
	Writer
}

type LocalFileReaderWriter struct{}

func (l *LocalFileReaderWriter) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (l *LocalFileReaderWriter) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

type S3FileReaderWriter struct {
	s3Client *s3.S3
	bucket   string
}

func (s *S3FileReaderWriter) Read(path string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	}

	result, err := s.s3Client.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3FileReaderWriter) Write(path string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   aws.ReadSeekCloser(io.NopCloser(bytes.NewReader(data))),
	}

	_, err := s.s3Client.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to put object to S3: %w", err)
	}

	return nil
}

type FileReaderWriterFactory struct {
	s3Client *s3.S3
	bucket   string
}

func NewFileReaderWriterFactory(s3Client *s3.S3, bucket string) *FileReaderWriterFactory {
	return &FileReaderWriterFactory{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

func (f *FileReaderWriterFactory) NewFileReaderWriter(isS3 bool) ReaderWriter {
	if isS3 {
		return &S3FileReaderWriter{
			s3Client: f.s3Client,
			bucket:   f.bucket,
		}
	}
	return &LocalFileReaderWriter{}
}
