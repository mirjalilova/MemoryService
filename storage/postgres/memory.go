package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	pb "github.com/mirjalilova/MemoryService/genproto/memory"
)

type MemoryRepo struct {
	db *sql.DB
}

func NewMemoryRepo(db *sql.DB) *MemoryRepo {
	return &MemoryRepo{db: db}
}

func (r *MemoryRepo) Create(req *pb.MemoryCreate) (*pb.Void, error) {
	res := &pb.Void{}

	location := fmt.Sprintf("(%.6f, %.6f)", req.Locations.Latitude, req.Locations.Longitude)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return res, fmt.Errorf("failed to parse date: %v", err)
	}

	tagsArray := pq.Array(req.Tags)

	query := `
        INSERT INTO memories (user_id, title, description, date, tags, location, place_name, type, privacy)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err = r.db.Exec(
		query,
		req.UserId,
		req.Title,
		req.Description,
		date,
		tagsArray,
		location,
		req.PlaceName,
		req.Type,
		req.Privacy,
	)

	if err != nil {
		return res, fmt.Errorf("failed to execute insert query: %v", err)
	}

	return res, nil
}

func (r *MemoryRepo) Get(id *pb.GetById) (*pb.MemoryRes, error) {
	res := &pb.MemoryRes{}

	query := `
        SELECT title, description, date, tags, location, place_name, type, privacy
        FROM memories
        WHERE id = $1 AND user_id = $2 AND deleted_at = 0
    `

	row := r.db.QueryRow(query, id.Id, id.UserId)

	var date time.Time
	var tags pq.StringArray
	var locationString string

	err := row.Scan(
		&res.Title,
		&res.Description,
		&date,
		&tags,
		&locationString,
		&res.PlaceName,
		&res.Type,
		&res.Privacy,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("memory not found or not accessible")
	} else if err != nil {
		return nil, fmt.Errorf("failed to scan row: %v", err)
	}

	res.Locations, err = parsePoint(locationString)
	if err != nil {
		return nil, fmt.Errorf("error parsing start point: %v", err)
	}

	res.Tags = tags
	res.Date = date.Format("2006-01-02")

	return res, nil
}

func (r *MemoryRepo) GetAll(req *pb.GetAllReq) (*pb.GetAllRes, error) {
	res := &pb.GetAllRes{}

	query := `
        SELECT title, description, date, tags, location::TEXT, place_name, type, privacy
        FROM memories
        WHERE user_id = $1 AND deleted_at = 0
    `

	var args []interface{}
	var conditions []string

	args = append(args, req.UserId)

	if req.EndDate != "" {
		args = append(args, req.EndDate)
		conditions = append(conditions, fmt.Sprintf("date <= $%d", len(args)))
	}

	if req.StartDate != "" {
		args = append(args, req.StartDate)
		conditions = append(conditions, fmt.Sprintf("date >= $%d", len(args)))
	}

	if req.Tag != "" {
		args = append(args, req.Tag)
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(tags)", len(args)))
	}

	if req.Type != "" {
		args = append(args, req.Type)
		conditions = append(conditions, fmt.Sprintf("type = $%d", len(args)))
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf(" AND %s", strings.Join(conditions, " AND "))
	}

	var defaultLimit int32
	err := r.db.QueryRow("SELECT COUNT(1) FROM memories WHERE deleted_at = 0 AND user_id = $1", req.UserId).Scan(&defaultLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %v", err)
	}

	if req.Filter.Limit == 0 {
		req.Filter.Limit = defaultLimit
	}

	args = append(args, req.Filter.Limit, req.Filter.Offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var memory pb.MemoryRes
		var date time.Time
		var tags pq.StringArray
		var locationStr string

		err := rows.Scan(
			&memory.Title,
			&memory.Description,
			&date,
			&tags,
			&locationStr,
			&memory.PlaceName,
			&memory.Type,
			&memory.Privacy,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		memory.Locations, err = parsePoint(locationStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing location point: %v", err)
		}

		memory.Tags = tags
		memory.Date = date.Format("2006-01-02")

		res.Memories = append(res.Memories, &memory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
	}
	res.Count = int32(len(res.Memories))

	return res, nil
}

func (r *MemoryRepo) Update(req *pb.MemoryUpdate) (*pb.Void, error) {
	res := &pb.Void{}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return res, fmt.Errorf("failed to parse date: %v", err)
	}

	location := fmt.Sprintf("(%.6f, %.6f)", req.Locations.Latitude, req.Locations.Longitude)

	query := `
        UPDATE memories SET updated_at = NOW()
    `
	var arg []interface{}
	var conditions []string

	if req.Title != "" && req.Title != "string" {
		arg = append(arg, req.Tags)
		conditions = append(conditions, fmt.Sprintf("title = $%d", len(arg)))
	}

	if req.Description != "" && req.Description != "string" {
		arg = append(arg, req.Description)
		conditions = append(conditions, fmt.Sprintf("description = $%d", len(arg)))
	}

	if req.Date != "" && req.Date != "string" {
		arg = append(arg, date)
		conditions = append(conditions, fmt.Sprintf("date = $%d", len(arg)))
	}

	if len(req.Tags) > 0 {
		tagsArray := pq.Array(req.Tags)
		arg = append(arg, tagsArray)
		conditions = append(conditions, fmt.Sprintf("tags = $%d", len(arg)))
	}

	if location != "" && location != "string" {
		arg = append(arg, location)
		conditions = append(conditions, fmt.Sprintf("location = $%d", len(arg)))
	}

	if req.PlaceName != "" && req.PlaceName != "string" {
		arg = append(arg, req.Date)
		conditions = append(conditions, fmt.Sprintf("place_name = $%d", len(arg)))
	}

	if req.Privacy != "" && req.Privacy != "string" {
		arg = append(arg, req.Privacy)
		conditions = append(conditions, fmt.Sprintf("privacy = $%d", len(arg)))
	}

	if len(conditions) > 0 {
		query += ", " + strings.Join(conditions, ", ")
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", len(arg)+1, len(arg)+2)
	arg = append(arg, req.Id, req.UserId)

	_, err = r.db.Exec(query, arg...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *MemoryRepo) Delete(id *pb.GetById) (*pb.Void, error) {
	res := &pb.Void{}

	query := `
        UPDATE memories SET deleted_at = NOW()
        WHERE id = $1 AND user_id = $2
    `

	_, err := r.db.Exec(query, id.Id, id.UserId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func parsePoint(pointStr string) (*pb.Point, error) {
	pointStr = strings.Trim(pointStr, "()")
	pointStr = strings.Replace(pointStr, ",", ", ", 1)
	var lat, lon float64
	_, err := fmt.Sscanf(pointStr, "%f, %f", &lat, &lon)
	if err != nil {
		return nil, fmt.Errorf("error parsing point: %v", err)
	}
	return &pb.Point{Latitude: lat, Longitude: lon}, nil
}

func (r *MemoryRepo) GetMemoriesOfOthers(user_id *pb.GetByUser) (*pb.GetAllRes, error) {
	res := &pb.GetAllRes{}

    query := `
        SELECT title, description, date, tags, location::TEXT, place_name, type, privacy
        FROM memories
        WHERE deleted_at = 0 AND id = ANY (
                                SELECT memory_id 
                                FROM shares 
								WHERE $1 = ANY(shared_with))
    `

    rows, err := r.db.Query(query, user_id.UserId)
    if err!= nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }

    defer rows.Close()

    for rows.Next() {
        var memory pb.MemoryRes
        var date time.Time
        var tags pq.StringArray
        var locationStr string

        err := rows.Scan(
            &memory.Title,
            &memory.Description,
            &date,
            &tags,
            &locationStr,
            &memory.PlaceName,
            &memory.Privacy,
			&memory.Type,
        )
		if err!= nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }

		memory.Locations, err = parsePoint(locationStr)
		if err!= nil {
            return nil, fmt.Errorf("error parsing location point: %v", err)
        }
		memory.Tags = tags
		memory.Date = date.Format("2006-01-02")
		res.Memories = append(res.Memories, &memory)

    }

    if err := rows.Err(); err!= nil {
        return nil, fmt.Errorf("failed to iterate rows: %v", err)
    }

    res.Count = int32(len(res.Memories))

    return res, nil
}
