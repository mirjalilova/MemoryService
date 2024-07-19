package service

import (
	"context"
	"log"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	st "github.com/mirjalilova/MemoryService/storage/postgres"
)

type MemoryService struct {
	storage st.Storage
	pb.UnimplementedMemoryServiceServer
}

func NewMemoryService(storage *st.Storage) *MemoryService {
	return &MemoryService{
		storage: *storage,
	}
}

func (s *MemoryService) Create(ctx context.Context, req *pb.MemoryCreate) (*pb.Void, error) {
	res, err := s.storage.MemoryS.Create(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *MemoryService) Get(ctx context.Context, req *pb.GetById) (*pb.MemoryRes, error) {
	res, err := s.storage.MemoryS.Get(req)
	if err != nil {
		log.Printf("Error getting memory: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *MemoryService) Update(ctx context.Context, req *pb.MemoryUpdate) (*pb.Void, error) {
	res, err := s.storage.MemoryS.Update(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *MemoryService) Delete(ctx context.Context, req *pb.GetById) (*pb.Void, error) {
	res, err := s.storage.MemoryS.Delete(req)
	if err != nil {
		log.Printf("Error deleting memory: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *MemoryService) GetAll(ctx context.Context, req *pb.GetAllReq) (*pb.GetAllRes, error) {
	res, err := s.storage.MemoryS.GetAll(req)
	if err != nil {
		log.Printf("Error evaluating memory: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *MemoryService) GetMemoriesOfOthers(ctx context.Context, req *pb.GetByUser) (*pb.GetAllRes, error) {
	res, err := s.storage.MemoryS.GetMemoriesOfOthers(req)
    if err != nil {
        log.Printf("Error getting memories of others: %v", err)
        return nil, err
    }
    return res, nil
}