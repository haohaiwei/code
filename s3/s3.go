package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/juju/ratelimit"
)


func main() {
	ak := flag.String("ak", "", "ak")
	sk := flag.String("sk", "", "sk")
	concurrent := flag.Int("concurrent", 1, "how many file concurrent upload")
	bucket := flag.String("bucket", "", "bucket")
	fsize := flag.Int64("fsize", 1, "fsize size")
	s3endpoint := flag.String("s3endponit", "", "s3endponit")
	region := flag.String("region","cn-east-1","region" )
	key := flag.String("key", "key", "filekey" )
	part_size := flag.Int64("part_size", 5, "part size")
	speed := flag.Int64("speed", 100, "MB")
	flag.Parse()
	b := float64(*speed)
	awsConfig := &aws.Config{
		Region:      aws.String(*region),
		Endpoint:    aws.String(*s3endpoint),
		Credentials: credentials.NewStaticCredentials(*ak, *sk, ""),
		S3ForcePathStyle: aws.Bool(true),
	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(awsConfig))

	// Create an uploader with the session and default options
	//uploader := s3manager.NewUploader(sess)

	// Create an uploader with the session and custom options
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = *part_size * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = *concurrent           // default is 5
	})

	//open the file
	f, err := os.Open("/dev/zero")
	if err != nil {
		fmt.Printf("failed to open: %v", err)
		return
	}
	defer f.Close()
	rdr := io.NewSectionReader(f, 0, *fsize * 1024 * 1024)
	limit := ratelimit.NewBucketWithRate(b * 1024 * 1024, *speed * 1024 * 1024)
	body := ratelimit.Reader(rdr, limit)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
		Body:   body,
	})

	//in case it fails to upload
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
}
