package service

import (
	"context"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
	st "github.com/mirjalilova/MemoryService/storage/postgres"
)

type CommentService struct {
	storage st.Storage
	pb.UnimplementedCommentServiceServer
}

func NewCommentService(storage *st.Storage) *CommentService {
	return &CommentService{
		storage: *storage,
	}
}

func (s *CommentService) Create(ctx context.Context, req *pb.CommentCreate) (*pb.Void, error) {
	res, err := s.storage.CommentS.Create(req)
	if err!= nil {
        return nil, err
    }
	return res, nil
}

func (s *CommentService) Get(ctx context.Context, req *pb.GetById) (*pb.Comment, error) {
	res, err := s.storage.CommentS.Get(req)
    if err!= nil {
        return nil, err
    }
    return res, nil
}

func (s *CommentService) Update(ctx context.Context, req *pb.CommentUpdate) (*pb.Void, error) {
	res, err := s.storage.CommentS.Update(req)
	if err!= nil {
        return nil, err
    }
	return res, nil
}

func (s *CommentService) Delete(ctx context.Context, req *pb.GetById) (*pb.Void, error) {
	res, err := s.storage.CommentS.Delete(req)
	if err!= nil {
        return nil, err
    }
	return res, nil
}
