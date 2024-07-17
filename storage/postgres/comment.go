package postgres

import (
	"database/sql"
	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(req *pb.CommentCreate) (*pb.Void, error) {
	return nil, nil
}

func (r *CommentRepo) Get(req *pb.GetById) (*pb.Comment, error) {
    return nil, nil
}

func (r *CommentRepo) Update(req *pb.GetById) (*pb.Void, error) {
    return nil, nil
}

func (r *CommentRepo) Delete(id *pb.GetById) (*pb.Void, error) {
    return nil, nil
}

