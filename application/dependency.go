package application

import (
	"context"
	"time"

	"github.com/Martin-Jast/go-microservice/persistence"
)

type IService interface {
	CreateBaseDocument(ctx context.Context, data string) (id string, err error)
	DeleteBaseDocument(ctx context.Context, id string) error
	GetBaseDocumentByID(ctx context.Context, id string) (*persistence.BaseModel, error)
	GetAllCreatedSince(ctx context.Context, date time.Time) ([]persistence.BaseModel, error)
}
