package postgres

import (
	"database/sql"
	"fmt"

	"github.com/mirjalilova/MemoryService/storage"

	"golang.org/x/exp/slog"

	_ "github.com/lib/pq"
	"github.com/mirjalilova/MemoryService/config"
)

type Storage struct {
	Db       *sql.DB
	CommentS storage.CommentI
	MediaS   storage.MediaI
	MemoryS  storage.MemoryI
	ShareS   storage.ShareI
}

func NewPostgresStorage(config config.Config) (*Storage, error) {
	conn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d sslmode=disable",
		config.DB_HOST,
		config.DB_USER,
		config.DB_NAME,
		config.DB_PASSWORD,
		config.DB_PORT,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	slog.Info("connected to db")

	return &Storage{
		Db:       db,
		CommentS: NewCommentRepo(db),
		MediaS:   NewMediaRepo(db),
		MemoryS:  NewMemoryRepo(db),
		ShareS:   NewShareRepo(db),
	}, nil
}
