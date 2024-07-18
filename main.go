package main

import (
	"log"
	"net"

	"golang.org/x/exp/slog"

	cf "github.com/mirjalilova/MemoryService/config"
	"github.com/mirjalilova/MemoryService/service"

	"github.com/mirjalilova/MemoryService/kafka"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	"github.com/mirjalilova/MemoryService/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	config := cf.Load()
	db, err := postgres.NewPostgresStorage(config)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return
	}
	defer db.Db.Close()

	memoryService := service.NewMemoryService(db)

	brokers := []string{"localhost:9092"}

	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "create-memory", "eval", kafka.MemoryCreateHandler(memoryService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-memory' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}

	if err := kcm.RegisterConsumer(brokers, "update-memory", "eval", kafka.MemoryUpdateHandler(memoryService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'update-memory' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}

	listener, err := net.Listen("tcp", config.MEMORY_PORT)
	if err != nil {
		slog.Error("can't listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterCommentServiceServer(s, service.NewCommentService(db))
	pb.RegisterMediaServiceServer(s, service.NewMediaService(db))
	pb.RegisterMemoryServiceServer(s, service.NewMemoryService(db))
	pb.RegisterShareServiceServer(s, service.NewShareService(db))

	slog.Info("server started port", config.MEMORY_PORT)
	if err := s.Serve(listener); err != nil {
		slog.Error("can't serve: %v", err)
		return
	}
}
