package persistence

import (
	"context"
	"time"
)

type SQLAdapter struct {
	sqlConnection *interface{}
}


func NewSQLAdapter(dbConnection *interface{}) SQLAdapter {
	return SQLAdapter{
		sqlConnection: dbConnection,
	}
}

func (s SQLAdapter) Create(ctx context.Context, document BaseModel) (id string, err error) {
	return "", nil
}

func (s SQLAdapter) GetByID(ctx context.Context, id string) (doc *BaseModel, err error) {
	return nil, nil
}

func (s SQLAdapter) Delete(ctx context.Context, id string) error {
	return nil
}

func (s SQLAdapter) GetAllCreatedSince(ctx context.Context, date time.Time) (docs []BaseModel, err error) {
	return nil, nil
}