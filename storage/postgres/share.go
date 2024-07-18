package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type ShareRepo struct {
	db *sql.DB
}

func NewShareRepo(db *sql.DB) *ShareRepo {
	return &ShareRepo{db: db}
}

func (r *ShareRepo) Share(req *pb.ShareCreate) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT user_id FROM memories WHERE id = $1`
	var userId string
	err := r.db.QueryRow(query, req.MemoryId).Scan(&userId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory not found")
	} else if err != nil {
		return nil, err
	}
	if userId != req.UserId {
		return nil, fmt.Errorf("not authorized to create media for this memory")
	}

	usersArray := pq.Array(req.SharedWith)

	query = `INSERT INTO shares (memory_id, shared_with) VALUES ($1, $2)`
	_, err = r.db.Exec(query, req.MemoryId, usersArray)
	if err!= nil {
		return nil, err
    }

	return res, nil
}

func (r *ShareRepo) Updateshare(req *pb.ShareDelete) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT user_id FROM memories WHERE id = $1`
	var userId string
	err := r.db.QueryRow(query, req.MemoryId).Scan(&userId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory not found")
	} else if err != nil {
		return nil, err
	}
	if userId != req.UserId {
		return nil, fmt.Errorf("not authorized to create media for this memory")
	}

	usersArray := pq.Array(req.UnsharedWith)

	query = `INSERT INTO shares (memory_id, shared_with) VALUES ($1, $2)`
	_, err = r.db.Exec(query, req.MemoryId, usersArray)
	if err!= nil {
		return nil, err
    }

	return res, nil
}

func (r *ShareRepo) Get(req *pb.GetById) (*pb.ShareRes, error) {
	res := &pb.ShareRes{}

	var shareWith pq.StringArray

	query := `SELECT s.shared_with, m.title 
				FROM shares s
				JOIN memories m ON s.memory_id = m.id
				WHERE s.memory_id = $1 AND m.user_id = $2`
		
	row := r.db.QueryRow(query, req.Id, req.UserId)	
	err := row.Scan(&shareWith, &res.MemonyTitle)
	if err == sql.ErrNoRows {
        return nil, fmt.Errorf("share not found")
    } else if err!= nil {
        return nil, err
    }

	res.SharedWith = shareWith
	
	return res, nil
}

