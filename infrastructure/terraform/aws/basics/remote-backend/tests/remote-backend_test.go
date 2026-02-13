package tests

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoteBackend(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	id := strings.ToLower(random.UniqueId())
	expectedEnvironment := "dev"
	expectedRegion := "ap-northeast-1"
	expectedBucketName := fmt.Sprintf("bucket-%s-%s", ts, id)
	expectedTableName := fmt.Sprintf("lock-%s-%s", ts, id)
	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]any{
			"environment":         expectedEnvironment,
			"aws_region":          expectedRegion,
			"bucket_name":         expectedBucketName,
			"dynamodb_table_name": expectedTableName,
		},
	})
	defer terraform.Destroy(t, opts)
	_, err := terraform.InitAndApplyE(t, opts)
	require.NoError(t, err)

	t.Run("outputs", func(t *testing.T) {
		tests := []struct {
			name       string
			outputName string
			want       string
		}{
			{
				name:       "aws region",
				outputName: "aws_region",
				want:       expectedRegion,
			},
			{
				name:       "s3 bucket name",
				outputName: "s3_bucket_name",
				want:       expectedBucketName,
			},
			{
				name:       "dynamodb table name",
				outputName: "dynamodb_table_name",
				want:       expectedTableName,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got := terraform.Output(t, opts, tt.outputName)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("s3 bucket configuration", func(t *testing.T) {
		bucketName := terraform.Output(t, opts, "s3_bucket_name")
		aws.AssertS3BucketExists(t, expectedRegion, bucketName)
		s3Client := aws.NewS3Client(t, expectedRegion)
		tests := []struct {
			name      string
			assertion func(t *testing.T)
		}{
			{
				name: "versioning enabled",
				assertion: func(t *testing.T) {
					got, err := s3Client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
						Bucket: awssdk.String(bucketName),
					})
					require.NoError(t, err)
					assert.Equal(t, s3types.BucketVersioningStatusEnabled, got.Status)
				},
			},
			{
				name: "encryption configured",
				assertion: func(t *testing.T) {
					got, err := s3Client.GetBucketEncryption(ctx, &s3.GetBucketEncryptionInput{
						Bucket: awssdk.String(bucketName),
					})
					require.NoError(t, err)
					require.NotEmpty(t, got.ServerSideEncryptionConfiguration.Rules)
					assert.Equal(t, s3types.ServerSideEncryptionAes256, got.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
				},
			},
			{
				name: "ssl enforcement policy",
				assertion: func(t *testing.T) {
					got, err := s3Client.GetBucketPolicy(ctx, &s3.GetBucketPolicyInput{
						Bucket: awssdk.String(bucketName),
					})
					require.NoError(t, err)
					assert.Contains(t, awssdk.ToString(got.Policy), "aws:SecureTransport")
				},
			},
			{
				name: "public access blocked",
				assertion: func(t *testing.T) {
					got, err := s3Client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
						Bucket: awssdk.String(bucketName),
					})
					require.NoError(t, err)
					config := got.PublicAccessBlockConfiguration
					require.NotNil(t, config)
					assert.True(t, awssdk.ToBool(config.BlockPublicAcls))
					assert.True(t, awssdk.ToBool(config.BlockPublicPolicy))
					assert.True(t, awssdk.ToBool(config.IgnorePublicAcls))
					assert.True(t, awssdk.ToBool(config.RestrictPublicBuckets))
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				tt.assertion(t)
			})
		}
	})

	t.Run("dynamodb table configuration", func(t *testing.T) {
		tableName := terraform.Output(t, opts, "dynamodb_table_name")
		dynamoDBClient := aws.NewDynamoDBClient(t, expectedRegion)
		tests := []struct {
			name      string
			assertion func(t *testing.T)
		}{
			{
				name: "table is active",
				assertion: func(t *testing.T) {
					got, err := dynamoDBClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
						TableName: awssdk.String(tableName),
					})
					require.NoError(t, err)
					assert.Equal(t, dynamodbtypes.TableStatusActive, got.Table.TableStatus)
				},
			},
			{
				name: "hash key configured correctly",
				assertion: func(t *testing.T) {
					got, err := dynamoDBClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
						TableName: awssdk.String(tableName),
					})
					require.NoError(t, err)
					hasCorrectHashKey := false
					for _, attr := range got.Table.AttributeDefinitions {
						if awssdk.ToString(attr.AttributeName) == "LockID" && attr.AttributeType == dynamodbtypes.ScalarAttributeTypeS {
							hasCorrectHashKey = true
							break
						}
					}
					assert.True(t, hasCorrectHashKey)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				tt.assertion(t)
			})
		}
	})
}
