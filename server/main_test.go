package server

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/Martin-Jast/go-microservice/application"
	"github.com/Martin-Jast/go-microservice/persistence"
	"github.com/Martin-Jast/go-microservice/utils"
	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"
)

type testHandler struct {
	dbAddapter persistence.PersistenceAdapter
	application application.IService
}

func createTestHandler() (*testHandler, context.Context) {
	ctx := context.Background()
	err := utils.SetupEnvVars("../.env")
	if err != nil {
		panic(err)
	}
	// Start by connecting to DB Clients
	mongoClient, err := persistence.CreateMongoConnection(ctx, os.Getenv("MONGO_STRING"))
	if err != nil {
		panic(err)
	}
	// Start Adapters
	mongoAdapter := persistence.NewMongoAdapter(mongoClient)

	// Start Application
	service := application.NewService(&mongoAdapter)

	th := new(testHandler)
	th.application = service
	th.dbAddapter = mongoAdapter

	return th, ctx
}

func (th testHandler) createHTTPExpect(t *testing.T) *httpexpect.Expect {
	service := th.application
	reqShutdown  := make(chan bool)
	handler := NewServer(service, reqShutdown, nil)

	api := mux.NewRouter()
	api.PathPrefix("/").Handler(handler)

	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(api),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	})
}
