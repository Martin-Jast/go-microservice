package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SQLAdapter struct {
	sqlConnection *sql.DB
}


func NewSQLAdapter(dbConnection *sql.DB) SQLAdapter {
	return SQLAdapter{
		sqlConnection: dbConnection,
	}
}

func (s SQLAdapter) Create(ctx context.Context, document BaseModel) (id string, err error) {
	lId := primitive.NewObjectID()
	_, err = s.sqlConnection.Exec("INSERT INTO test (id, data, created_at, deleted_at) VALUES (?, ?, ?)", lId.Hex(), document.Data, document.CreatedAt, document.DeletedAt)
	if err != nil {
		return "", err
	}
    return lId.Hex(), nil
}

func (s SQLAdapter) GetByID(ctx context.Context, id string) (doc *BaseModel, err error) {
	row := s.sqlConnection.QueryRow("SELECT * FROM album WHERE id = ?", id)
	elem := BaseModel{}
    if err := row.Scan(&elem.ID, &elem.Data, &elem.CreatedAt, &elem.DeletedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no document with id: %s", id)
		}
    	return nil, err
    }
    return &elem, nil
}

func (s SQLAdapter) Delete(ctx context.Context, id string) error {
	_, err := s.sqlConnection.Exec("DELETE FROM test WHERE id = ?", id)
	if err != nil {
		return err
	}
    return nil
}

func (s SQLAdapter) GetAllCreatedSince(ctx context.Context, date time.Time) (docs []BaseModel, err error) {
	rows, err := s.sqlConnection.Query("SELECT * FROM album WHERE created_at > ?", date)
	if err != nil {
        return nil, err
    }
	defer rows.Close()
	result := []BaseModel{}
	for rows.Next() {
		elem := BaseModel{}
		if err := rows.Scan(&elem.ID, &elem.Data, &elem.CreatedAt, &elem.DeletedAt); err != nil {
			return nil, err
		}
		result = append(result, elem)
    }
	if err := rows.Err(); err != nil {
        return nil, err
    }
    return result, nil
}
