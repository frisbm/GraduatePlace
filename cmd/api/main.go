package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frisbm/graduateplace/pkg/tasks"

	"github.com/hibiken/asynq"

	"github.com/frisbm/graduateplace/pkg/stack/document"

	"github.com/frisbm/graduateplace/pkg/services/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pressly/goose/v3"

	"github.com/frisbm/graduateplace/pkg/store"

	"github.com/go-chi/jwtauth/v5"

	"github.com/frisbm/graduateplace/pkg/services/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"

	"github.com/frisbm/graduateplace/pkg/config"
	"github.com/frisbm/graduateplace/pkg/stack/user"
)

const (
	API_PATH = "/api"
	TIMEOUT  = 10 * time.Second
)

type Route interface {
	Public(c chi.Router)
	Private(c chi.Router)
	Admin(c chi.Router)
}

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())

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

	if err = goose.Up(database, "pkg/store/database/migrations"); err != nil {
		log.Fatalf("failed migrating with goose: %v", err)
	}
	db, err := store.Prepare(ctx, database)
	if err != nil {
		log.Fatalf("failed preparing preparing statments: %v", err)
	}
	defer db.Close()

	// Asynq,Redis, & Task Manager Setup
	// ###################################################
	asynqClient := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%v:%v", config.RedisHost, config.RedisPort),
			Password: config.RedisPassword,
		},
	)
	defer asynqClient.Close()
	taskManager := tasks.NewTaskManager(asynqClient)

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

	// Auth Service
	// ###################################################
	authService := auth.NewAuthService(config.JWTSecretKey)

	// Create and Initialize "Stacks"
	// Stacks tightly bind routes, manager, and stores
	// ###################################################
	userStack := user.NewStack(db, s3, taskManager, authService)
	documentStack := document.NewStack(db, s3, taskManager)

	// Router Setup
	// ###################################################
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(TIMEOUT))

	public := r.Group(nil)

	authMiddleware := auth.NewAuthMiddleware(db, authService)

	private := r.Group(nil)
	private.Use(jwtauth.Verifier(authService.GetTokenAuth()))
	private.Use(authMiddleware.Private)

	admin := r.Group(nil)
	admin.Use(jwtauth.Verifier(authService.GetTokenAuth()))
	admin.Use(authMiddleware.Admin)

	routes := []Route{
		userStack.Router,
		documentStack.Router,
	}

	for _, route := range routes {
		route.Public(public)
		route.Private(private)
		route.Admin(admin)
	}

	router := chi.NewRouter()
	router.Mount(API_PATH, r)

	// Start server with graceful shutdown
	// ###################################################
	server := &http.Server{Addr: ":8080", Handler: router}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Println("Hi")
		shutdownCtx, _ := context.WithTimeout(ctx, TIMEOUT)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("forcing server shutdown due to graceful timeout...")
			}
		}()

		err = server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("shutting down server...")
		ctxCancel()
	}()

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-ctx.Done()
}
