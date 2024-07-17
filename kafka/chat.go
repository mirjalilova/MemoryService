package kafka

import (
	// "context"
	// "encoding/json"
	// "log"

	// pb "github.com/mirjalilova/MemoryService/genproto/memory"
	// "github.com/mirjalilova/MemoryService/service"
)

// func ChatCreateHandler(ch *service.ChatService) func(message []byte) {
// 	return func(message []byte) {
// 		var chat pb.ChatCreate
// 		if err := json.Unmarshal(message, &chat); err != nil {
// 			log.Printf("Cannot unmarshal JSON: %v", err)
// 			return
// 		}

// 		_, err := ch.Create(context.Background(), &chat)
// 		if err != nil {
// 			log.Printf("Cannot create chat via Kafka: %v", err)
// 			return
// 		}
// 		log.Printf("Created chat")
// 	}
// }
