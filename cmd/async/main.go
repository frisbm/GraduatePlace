package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/frisbm/graduateplace/pkg/services/s3"
	"github.com/frisbm/graduateplace/pkg/store/document"
	"github.com/frisbm/graduateplace/pkg/tasks/handlers"

	"github.com/frisbm/graduateplace/pkg/tasks"

	"github.com/frisbm/graduateplace/pkg/constants"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"

	"github.com/frisbm/graduateplace/pkg/config"
	"github.com/frisbm/graduateplace/pkg/store"
)

func main() {
	// Config
	// ###################################################
	config := config.NewConfig("config.json")

	// Database Connection & Migration
	// ###################################################
	connectionString := fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBSSLMode,
	)
	// Open postgres
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer database.Close()

	store := store.New(database)
	dbStore := document.NewStore(store)

	// AWS Setup
	// ###################################################
	defaultRegion := "us-east-1"
	awsEndpoint := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			// If ENVIRONMENT is local
			if config.Environment == "local" {
				return aws.Endpoint{
					PartitionID:       "aws",
					URL:               config.AwsEndpoint,
					SigningRegion:     defaultRegion,
					HostnameImmutable: true,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)
	awsCredentials := aws.CredentialsProviderFunc(
		func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     config.AwsAccessKeyId,
				SecretAccessKey: config.AwsSecretAccessKey,
			}, nil
		},
	)

	cfg := aws.Config{
		Region:                      defaultRegion,
		Credentials:                 awsCredentials,
		EndpointResolverWithOptions: awsEndpoint,
	}

	s3 := s3.NewS3(cfg)

	// Asynq & Redis Setup
	// ###################################################
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%v:%v", config.RedisHost, config.RedisPort),
			Password: config.RedisPassword,
		},
		asynq.Config{
			Concurrency: 4,
			Queues: map[string]int{
				constants.HIGH_PRIORITY_QUEUE: 3,
				constants.LOW_PRIORITY_QUEUE:  1,
			},
		},
	)
	mux := asynq.NewServeMux()
	mux.Handle(tasks.ProcessDocumentTask, handlers.NewDocumentProcessor(dbStore, s3))

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
