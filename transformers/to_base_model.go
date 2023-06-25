package transformers

import (
	"time"

	"github.com/Martin-Jast/go-microservice/persistence"
	"github.com/Martin-Jast/go-microservice/utils"
)

type BaseModelResponse struct {
	ID        *string `json:"id"`
	Data      string `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func ToBaseModelResponse(b persistence.BaseModel) BaseModelResponse {
	return BaseModelResponse{
		ID: utils.StrPnt(*b.ID),
		Data: b.Data,
		CreatedAt: b.CreatedAt,
		DeletedAt: b.DeletedAt,
	}
}

func ToBaseModelResponseArray(bs []persistence.BaseModel) []BaseModelResponse{
	response := make([]BaseModelResponse, len(bs))
	for i := range bs {
		response[i] = ToBaseModelResponse(bs[i])
	}
	return response
}