package kafka

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	"github.com/mirjalilova/MemoryService/service"
)

func MemoryCreateHandler(ch *service.MemoryService) func(message []byte) {
	return func(message []byte) {
		var memory pb.MemoryCreate
		if err := json.Unmarshal(message, &memory); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		_, err := ch.Create(context.Background(), &memory)
		if err != nil {
			log.Printf("Cannot create memory via Kafka: %v", err)
			return
		}
		log.Printf("Created memory")
	}
}

func MemoryUpdateHandler(ch *service.MemoryService) func(message []byte) {
	return func(message []byte) {
		var memory pb.MemoryUpdate
		if err := json.Unmarshal(message, &memory); err != nil {
			log.Printf("Cannot unmarshal JSON: %v", err)
			return
		}

		_, err := ch.Update(context.Background(), &memory)
		if err != nil {
			log.Printf("Cannot update memory via Kafka: %v", err)
			return
		}
		log.Printf("Updated memory")
	}
}

