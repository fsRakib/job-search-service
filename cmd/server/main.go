package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"job-search-service/internal/elastic"
	grpcHandler "job-search-service/internal/grpc"
	"job-search-service/internal/repository"
	"job-search-service/internal/service"
	pb "job-search-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Elasticsearch struct {
		URL   string `yaml:"url"`
		Index string `yaml:"index"`
	} `yaml:"elasticsearch"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return &config, nil
}

func main() {
	log.Println("Starting Job Search Service...")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	esClient, err := elastic.NewClient(config.Elasticsearch.URL)
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	ctx := context.Background()
	if err := esClient.CreateIndex(ctx, config.Elasticsearch.Index); err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}

	jobRepo := repository.NewJobRepository(esClient.ES, config.Elasticsearch.Index)
	jobService := service.NewJobService(jobRepo)
	jobHandler := grpcHandler.NewJobHandler(jobService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterJobServiceServer(grpcServer, jobHandler)

	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %d", config.Server.Port)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
