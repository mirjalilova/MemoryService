package service

import (
	"context"
	"log"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	st "github.com/mirjalilova/MemoryService/storage/postgres"
)

type MediaService struct {
	storage st.Storage
	pb.UnimplementedMediaServiceServer
}

func NewMediaService(storage *st.Storage) *MediaService {
	return &MediaService{
		storage: *storage,
	}
}

func (s *MediaService) Create(ctx context.Context, req *pb.MediaCreate) (*pb.Void, error) {
	res, err := s.storage.MediaS.Create(req)
	if err != nil {
		log.Printf("Error creating media: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *MediaService) Get(ctx context.Context, req *pb.GetById) (*pb.MediaRes, error) {
	res, err := s.storage.MediaS.Get(req)
	if err != nil {
		log.Printf("Error getting media: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *MediaService) Delete(ctx context.Context, req *pb.GetById) (*pb.Void, error) {
	res, err := s.storage.MediaS.Delete(req)
	if err != nil {
		log.Printf("Error deleting media: %v", err)
		return nil, err
	}
	return res, nil
}
