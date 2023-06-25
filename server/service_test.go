package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Martin-Jast/go-microservice/persistence"
	"github.com/Martin-Jast/go-microservice/transformers"
	"github.com/gavv/httpexpect"
	"github.com/stretchr/testify/assert"
)

func TestService_Mongo(t *testing.T) {
	tt := []BaseHandlerTest{
		// Create
		{
			Name:         "fail- malformed request",
			HTTPMethod:   "POST",
			Path:         "/base/create",
			Req:          "",
			ExpectedCode: http.StatusBadRequest,
		},
		{
			Name:         "success - create new document",
			HTTPMethod:   "POST",
			Path:         "/base/create",
			Req:          createBaseDocumentRequest{Data: "test-data"},
			ExpectedCode: http.StatusOK,
			AssertPosTestDBStates: func(ctx context.Context, th *testHandler, sr *SetupResult, res *httpexpect.Response) error {
				cr := transformers.CreateBaseDocResponse{}
				err := json.Unmarshal([]byte(res.Body().Raw()), &cr)
				if err != nil && err != io.EOF {
					return fmt.Errorf("invalid response")
				}
				doc, err := th.dbAddapter.GetByID(ctx, cr.ID)
				if err != nil {
					return err
				}
				
				assert.Equal(t, "test-data", doc.Data)
				return nil
			},
		},
		// Delete
		{
			Name:         "success - delete document",
			HTTPMethod:   "GET",
			SetupPreTestDBs: func(ctx context.Context, th *testHandler) *SetupResult {
				// Create a document in the DB
				doc1 := persistence.BaseModel{
					Data: "should be deleted",
				}
				res1, err := th.dbAddapter.Create(ctx, doc1)
				if err != nil {
					panic(err)
				}
				doc1.ID = &res1

				doc2 := persistence.BaseModel{
					Data: "should NOT be deleted",
				}
				res2, err := th.dbAddapter.Create(ctx, doc2)
				if err != nil {
					panic(err)
				}
				doc2.ID = &res2
				return &SetupResult{
					BaseDoc: []persistence.BaseModel{
						doc1,
						doc2,
					},
				}
			},
			MountPath: func(sr *SetupResult) (string, interface{}) {
				return fmt.Sprintf("/base/delete/%s", *sr.BaseDoc[0].ID), nil
			},
			Req:          createBaseDocumentRequest{Data: "test-data"},
			ExpectedCode: http.StatusOK,
			AssertPosTestDBStates: func(ctx context.Context, th *testHandler, sr *SetupResult, res *httpexpect.Response) error {
				id1 := sr.BaseDoc[0].ID
				id2 := sr.BaseDoc[1].ID
				_, err := th.dbAddapter.GetByID(ctx, *id1)
				assert.Error(t, err)
				doc2, err := th.dbAddapter.GetByID(ctx, *id2)
				if err != nil {
					// Cannot retrieve the deleted document
					return err
				}
				assert.Equal(t, "should NOT be deleted", doc2.Data)
				return nil
			},
		},
		// Get By ID
		{
			Name:         "success - Get document",
			HTTPMethod:   "GET",
			SetupPreTestDBs: func(ctx context.Context, th *testHandler) *SetupResult {
				// Create a document in the DB
				doc1 := persistence.BaseModel{
					Data: "test-data",
				}
				res1, err := th.dbAddapter.Create(ctx, doc1)
				if err != nil {
					panic(err)
				}
				doc1.ID = &res1
				return &SetupResult{
					BaseDoc: []persistence.BaseModel{
						doc1,
					},
				}
			},
			MountPath: func(sr *SetupResult) (string, interface{}) {
				return fmt.Sprintf("/base/%s", *sr.BaseDoc[0].ID), nil
			},
			Req:          createBaseDocumentRequest{Data: "test-data"},
			ExpectedCode: http.StatusOK,
			AssertPosTestDBStates: func(ctx context.Context, th *testHandler, sr *SetupResult, res *httpexpect.Response) error {
				id1 := sr.BaseDoc[0].ID
				doc1, err := th.dbAddapter.GetByID(ctx, *id1)
				assert.NoError(t, err)
				assert.Equal(t, "test-data", doc1.Data)
				return nil
			},
		},
		// Get All since
		{
			Name:         "success - Get all documents since 1 hour ago",
			HTTPMethod:   "GET",
			SetupPreTestDBs: func(ctx context.Context, th *testHandler) *SetupResult {
				// Clear database
				_ = th.dbAddapter.DeleteAll(ctx)
				// Create a document in the DB
				docs := []persistence.BaseModel{
				}
				for i:=0; i<10; i++ {
					t := time.Now().Add(-time.Minute*time.Duration(10*i))
					doc1 := persistence.BaseModel{
						Data: fmt.Sprintf("test-data-%d", i),
						CreatedAt: &t,
					}
					res1, err := th.dbAddapter.Create(ctx, doc1)
					if err != nil {
						panic(err)
					}
					doc1.ID = &res1
					docs = append(docs, doc1)
				}
				return &SetupResult{
					BaseDoc: docs,
				}
			},
			MountPath: func(sr *SetupResult) (string, interface{}) {
				return fmt.Sprintf("/base/since/%s", time.Now().Add(-time.Hour).UTC().Format("2006-01-02T15:04:05Z")), nil
			},
			ExpectedCode: http.StatusOK,
			AssertPosTestDBStates: func(ctx context.Context, th *testHandler, sr *SetupResult, res *httpexpect.Response) error {
				cr := []transformers.BaseModelResponse{}
				err := json.Unmarshal([]byte(res.Body().Raw()), &cr)
				if err != nil && err != io.EOF {
					return fmt.Errorf("invalid response")
				}
				assert.Equal(t, 7, len(cr))
				return nil
			},
		},
	}
	ExecHandlerTest(tt, t)
}