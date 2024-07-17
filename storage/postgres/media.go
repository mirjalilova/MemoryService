package postgres

import (
	"database/sql"
	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type MediaRepo struct {
	db *sql.DB
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{db: db}
}

func (r *MediaRepo) Create(req *pb.MediaCreate) (*pb.Void, error) {
	return nil, nil
}

func (r *MediaRepo) Get(req *pb.GetById) (*pb.Media, error) {
    return nil, nil
}

func (r *MediaRepo) Delete(req *pb.GetById) (*pb.Void, error) {
	return nil, nil
}