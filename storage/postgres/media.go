package postgres

import (
	"database/sql"
	"fmt"

	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type MediaRepo struct {
	db *sql.DB
}

func NewMediaRepo(db *sql.DB) *MediaRepo {
	return &MediaRepo{db: db}
}

func (r *MediaRepo) Create(req *pb.MediaCreate) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT user_id FROM memories WHERE id = $1`
	var userId string
	err := r.db.QueryRow(query, req.MemoryId).Scan(&userId)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory not found")
	} else if err != nil {
		return nil, err
	}
	if userId != req.UserId{
        return nil, fmt.Errorf("not authorized to create media for this memory")
    }

	query = `INSERT INTO media (memory_id, type, url) VALUES ($1, $2, $3)`

	_, err = r.db.Exec(query, req.MemoryId, req.Type, req.Url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MediaRepo) Get(req *pb.GetById) (*pb.MediaRes, error) {
	res := &pb.MediaRes{}

	query := `SELECT md.id, md.memory_id, md.type, md.url, md.created_at
				FROM media md 
				JOIN memories mr ON md.memory_id = mr.id
				WHERE md.memory_id = $1 AND mr.user_id = $2 AND md.deleted_at = 0`

	rows, err := r.db.Query(query, req.Id, req.UserId)
	if err!= nil {
        return nil, err
    }
	defer rows.Close()

	var createdAt string

	for rows.Next() {
		var media pb.Media
		err := rows.Scan(&media.Id, &media.MemoryId, &media.Type, &media.Url, &createdAt)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("media not found")
		} else if err!= nil {
			return nil, err
		}
		media.CreatedAt = createdAt

		res.Media = append(res.Media, &media)
	}
	
	return res, nil
}

func (r *MediaRepo) Delete(req *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

	query := `SELECT memory_id FROM media WHERE id = $1`
	var memoryId string
	err := r.db.QueryRow(query, req.Id).Scan(&memoryId)
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

	if userId!= req.UserId {
        return nil, fmt.Errorf("not authorized to delete media for this memory")
    }
	query = `UPDATE media SET deleted_at=EXTRACT(EPOCH FROM NOW()) WHERE id = $1`
	_, err = r.db.Exec(query, req.Id)
	if err!= nil {
        return nil, err
    }

	return res, nil
}
