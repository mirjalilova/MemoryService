package postgres

import (
	"database/sql"
	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type ShareRepo struct {
	db *sql.DB
}

func NewShareRepo(db *sql.DB) *ShareRepo {
	return &ShareRepo{db: db}
}

func (r *ShareRepo) Share(req *pb.ShareCreate) (*pb.Void, error) {
	return nil, nil
}

func (r *ShareRepo) Unshare(req *pb.ShareDelete) (*pb.Void, error) {
    return nil, nil
}

func (r *ShareRepo) Get(req *pb.GetById) (*pb.ShareRes, error) {
    return nil, nil
}

func (r *ShareRepo) Update(req *pb.GetById) (*pb.Void, error) {
    return nil, nil
}