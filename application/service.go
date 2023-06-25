package application

import (
	"context"
	"time"

	"github.com/Martin-Jast/go-microservice/persistence"
)

type Service struct {
	PersistenceAdapter persistence.PersistenceAdapter
}

func NewService(adpt persistence.PersistenceAdapter) Service {
	return Service{adpt}
}

func (s Service) CreateBaseDocument(ctx context.Context, data string) (id string, err error){
	return s.PersistenceAdapter.Create(ctx, persistence.BaseModel{
		Data: data,
	})
}

func (s Service) DeleteBaseDocument(ctx context.Context, id string) error {
	return s.PersistenceAdapter.Delete(ctx, id)
}


func (s Service) GetBaseDocumentByID(ctx context.Context, id string) (*persistence.BaseModel, error) {
	return s.PersistenceAdapter.GetByID(ctx, id)
}

func (s Service) GetAllCreatedSince(ctx context.Context, date time.Time) ([]persistence.BaseModel, error) {
	return s.PersistenceAdapter.GetAllCreatedSince(ctx, date)
}