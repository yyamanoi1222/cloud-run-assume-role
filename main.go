package main

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var roleArn = os.Getenv("ROLE_ARN")

type gcp struct{}

func (*gcp) GetIdentityToken() ([]byte, error) {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return []byte(""), err
	}

	ts, err := idtoken.NewTokenSource(ctx, "test-aud", option.WithCredentials(credentials))
	if err != nil {
		return []byte(""), err
	}
	t, err := ts.Token()
	if err != nil {
		return []byte(""), err
	}
	return []byte(t.AccessToken), nil
}

func getS3Client(ctx context.Context, g *gcp) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.Region = "ap-northeast-1"
	if err != nil {
		return nil, err
	}
	creds := stscreds.NewWebIdentityRoleProvider(sts.NewFromConfig(cfg), roleArn, g)
	cfg.Credentials = aws.NewCredentialsCache(creds)

	return s3.NewFromConfig(cfg), nil
}

func main() {
	ctx := context.Background()
	g := &gcp{}
	cl, err := getS3Client(ctx, g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}
	out, err := cl.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}
	for i := range out.Buckets {
		fmt.Printf("bucket: %s\n", *out.Buckets[i].Name)
	}
}
