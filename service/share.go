package service

import (
	"context"
	"log"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	st "github.com/mirjalilova/MemoryService/storage/postgres"
)

type ShareService struct {
	storage st.Storage
	pb.UnimplementedShareServiceServer
}

func NewShareService(storage *st.Storage) *ShareService {
	return &ShareService{
		storage: *storage,
	}
}

func (s *ShareService) Share(ctx context.Context, req *pb.ShareCreate) (*pb.Void, error) {
	res, err := s.storage.ShareS.Share(req)
	if err != nil {
		log.Printf("Error sharing share: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *ShareService) Get(ctx context.Context, req *pb.GetById) (*pb.ShareRes, error) {
	res, err := s.storage.ShareS.Get(req)
	if err != nil {
		log.Printf("Error getting share: %v", err)
		return nil, err
	}
	return res, nil
}

func (s *ShareService) Updateshare(ctx context.Context, req *pb.ShareDelete) (*pb.Void, error) {
	res, err := s.storage.ShareS.Updateshare(req)
	if err != nil {
		log.Printf("Error unsharing share: %v", err)
		return nil, err
	}
	return res, nil
}
