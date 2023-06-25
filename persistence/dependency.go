package persistence

import (
	"context"
	"time"
)

// BaseModel a simple generic DB document model to be used on this example
type BaseModel struct {
	ID *string `bson:"_id"`
	Data string
	CreatedAt *time.Time `bson:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}


// PersistenceAdapter defines how the application can communicate with a persistence Layer with no knowledge about how it is built
type PersistenceAdapter interface {
	Create(ctx context.Context, document BaseModel) (id string, err error)
	GetByID(ctx context.Context, id string) ( doc *BaseModel, err error)
	Delete(ctx context.Context, id string) error
	GetAllCreatedSince(ctx context.Context, date time.Time) (docs []BaseModel, err error)
	DeleteAll(ctx context.Context) error
}
