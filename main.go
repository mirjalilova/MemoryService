package main

import (
	"log"
	"net"

	"golang.org/x/exp/slog"

	cf "github.com/mirjalilova/MemoryService/config"
	"github.com/mirjalilova/MemoryService/kafka"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	"github.com/mirjalilova/MemoryService/service"
	"github.com/mirjalilova/MemoryService/storage/mongo"
	"github.com/mirjalilova/MemoryService/storage/postgres"

	"path/filepath"
	"runtime"

	"google.golang.org/grpc"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func main() {
	config := cf.Load()
	db, err := postgres.NewPostgresStorage(config)
	if err != nil {
		slog.Error("can't connect to db: %v", err)
		return
	}
	defer db.Db.Close()

	mdb, err := mongo.ConnectMongo()
	if err != nil {
		slog.Error("can't connect to mongo: %v", "err", err)
		return
	}

	evalService := service.NewEvaluationService(db)
	chatService := service.NewChatService(mdb)

	brokers := []string{"localhost:9092"}

	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "create-eval", "eval", kafka.EvaluationHandler(evalService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'create-eval' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "cre-chat", "chat", kafka.ChatCreateHandler(chatService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'cre-chat' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "upd-chat", "chat", kafka.ChatUpdateHandler(chatService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'upd-chat' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "cre-mess", "chat", kafka.MessageCreateHandler(chatService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'cre-mess' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}
	if err := kcm.RegisterConsumer(brokers, "upd-mess", "chat", kafka.MessageUpdateHandler(chatService)); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			log.Printf("Consumer for topic 'upd-mess' already exists")
		} else {
			log.Fatalf("Error registering consumer: %v", err)
		}
	}

	listener, err := net.Listen("tcp", ":7060")
	if err != nil {
		slog.Error("can't listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterCompanyServiceServer(s, service.NewCompanyService(db))
	pb.RegisterPositionServiceServer(s, service.NewPositionService(db))
	pb.RegisterDepartmentServiceServer(s, service.NewDepartmentService(db))
	pb.RegisterEvaluationServiceServer(s, service.NewEvaluationService(db))
	pb.RegisterChatServiceServer(s, service.NewChatService(mdb))
	pb.RegisterNotificationServiceServer(s, service.NewNotificationService(mdb))
	pb.RegisterGuideServiceServer(s, service.NewGuideService(db))

	slog.Info("server started port", config.COMPANY_PORT)
	if err := s.Serve(listener); err != nil {
		slog.Error("can't serve: %v", err)
		return
	}
}
