package postgres

import (
	"database/sql"
	"fmt"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) Create(req *pb.CommentCreate) (*pb.Void, error) {
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

    query = `INSERT INTO comments (memory_id, content, writer_id) VALUES ($1, $2, $3)`
    _, err = r.db.Exec(query, req.MemoryId, req.Content, req.UserId)
    if err != nil {
        return nil, err
    }

    return res, nil
}


func (r *CommentRepo) Get(req *pb.GetById) (*pb.Comment, error) {
	res := &pb.Comment{}

	query := `SELECT c.id, c.memory_id, c.writer_id, c.content, c.created_at, c.updated_at
	    FROM comments c 
		JOIN memories m ON c.memory_id = m.id
		WHERE c.memory_id = $1 AND m.user_id = $2 AND c.deleted_at = 0`

	row := r.db.QueryRow(query, req.Id, req.UserId)
	var createdAt, updatedAt string
	err := row.Scan(&res.Id, &res.MemoryId, &res.WriterId, &res.Content, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
        return nil, fmt.Errorf("comment not found")
    } else if err!= nil {
        return nil, err
    }

	res.CreatedAt = createdAt
	res.UpdatedAt = updatedAt

	return res, nil
}

func (r *CommentRepo) Update(req *pb.CommentUpdate) (*pb.Void, error) {
	res := &pb.Void{}

	query := `UPDATE comments SET content=$1, updated_at=NOW() WHERE id = $2`

	_, err := r.db.Exec(query, req.Content, req.Id)
	if err!= nil {
        return nil, err
    }
	
	return res, nil
}

func (r *CommentRepo) Delete(id *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT memory_id FROM comments WHERE id = $1`
	var memoryId string
	err := r.db.QueryRow(query, id.Id).Scan(&memoryId)
	if err == sql.ErrNoRows {
        return nil, fmt.Errorf("media not found")
    } else if err!= nil {
        return nil, err
    }

	query = `SELECT user_id FROM memories WHERE id = $1`
	var userId string
	err = r.db.QueryRow(query, memoryId).Scan(&userId)
	if err == sql.ErrNoRows {
        return nil, fmt.Errorf("memory not found")
    } else if err!= nil {
        return nil, err
    }

	if userId!= id.UserId {
        return nil, fmt.Errorf("not authorized to delete media for this memory")
    }
	query = `UPDATE comments SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id = $1`
	_, err = r.db.Exec(query, id.Id)
	if err!= nil {
        return nil, err
    }

	return res, nil
}
