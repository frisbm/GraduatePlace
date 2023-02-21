package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MatthewFrisby/thesis-pieces/pkg/constants"

	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"

	"github.com/MatthewFrisby/thesis-pieces/pkg/config"
	"github.com/MatthewFrisby/thesis-pieces/pkg/store"
	"github.com/MatthewFrisby/thesis-pieces/pkg/tasks"
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

	_ = store.New(database)

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
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleSendUserEmailTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
