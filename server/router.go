package server

import (
	"fmt"
	"net/http"

	"github.com/Martin-Jast/go-microservice/application"
	"github.com/gorilla/mux"
)

// New creates a new router
func NewServer(service application.IService, reqShutdown chan bool, middleware func(http.Handler) http.Handler) *mux.Router {
	router := mux.NewRouter()
	// router.HandleFunc("/heartbeat", h.HealthzHandler)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {fmt.Println("arrived")})

	// In case we want to add any root middlewares
	if middleware != nil {
		router.Use(middleware)
	}

	// endpoint to handle shutdown
	router.HandleFunc("/shutdown",func (w http.ResponseWriter, r *http.Request) { reqShutdown <-true})
	// Declare the prefix for which this service will handle request and assign it
	router.PathPrefix("/base").Handler(newServicePort(service))

	return router
}

type servicePort struct {
	*mux.Router
	service application.IService
}

func newServicePort(service application.IService) servicePort {
router := mux.NewRouter().PathPrefix("/base").Subrouter()
handler := servicePort{
	router,
	service,
}

router.Path("/create").
	Methods(http.MethodPost).HandlerFunc(handler.handleCreate)
router.Path("/delete/{id}").
	Methods(http.MethodGet).HandlerFunc(handler.handleDelete)
// Here we have a api to get documents since a date. 
// Another ( in my opinion better ) option would be to have a listAll api to get the documents and pass as queries the filters needed
router.Path("/since/{date}").
	Methods(http.MethodGet).HandlerFunc(handler.handleGetSince)
router.Path("/{id}").
	Methods(http.MethodGet).HandlerFunc(handler.handleGet)
	

return handler
}
