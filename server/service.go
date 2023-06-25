package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Martin-Jast/go-microservice/transformers"
	"github.com/Martin-Jast/go-microservice/utils"
	"github.com/gorilla/mux"
)

// Since there might be several endpoints it is better if we add those handler functions on separated files or a folder for each subrouter
// Here since it is not the case we are going to work with a single file

type createBaseDocumentRequest struct{
	Data string
}

// Build parses the request into our internal structure
func (cr *createBaseDocumentRequest) Build(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&cr)
	if err != nil && err != io.EOF {
		return fmt.Errorf("invalid request")
	}


	return nil
}

// Validate validates the request, should only check contract errors, never business logic
func (cr *createBaseDocumentRequest) Validate() error {
	var missingParams []string
	if cr.Data == "" {
		missingParams = append(missingParams, "Data")
	}
	if len(missingParams) > 0 {
		return fmt.Errorf("missing parameters: %s", strings.Join(missingParams, "; "))
	}
	return nil
}



// handleCreate handles the request for the creation of new documents
func (h servicePort) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Parse and Validate request
	req := &createBaseDocumentRequest{}
	req.Build(r)
	if err:= req.Validate(); err != nil {
		utils.WriteJson(err, w, 400)
		return;
	}
	// Deal with the request in application layer
	response, err := h.service.CreateBaseDocument(r.Context(), req.Data)
	if err != nil {
		utils.WriteError(fmt.Errorf("could not create document: %v", err), w, 500)
		return;
	}

	utils.WriteJson(transformers.CreateBaseDocResponse{ID: response}, w, 200)
}

// handleDelete handles the request for the deletion of new documents
func (h servicePort) handleDelete(w http.ResponseWriter, r *http.Request) {
	// Parse and Validate request
	id := mux.Vars(r)["id"]
	if id=="" {
		utils.WriteError(fmt.Errorf("missing id to delete"), w, 400)
		return;
	}
	// Deal with the request in application layer
	err := h.service.DeleteBaseDocument(r.Context(), id)
	if err != nil {
		utils.WriteError(fmt.Errorf("could not delete document: %v", err), w, 500)
		return;
	}

	utils.WriteJson(nil, w, 200)
}

// handleGet handles the request for getting documents by their ids
func (h servicePort) handleGet(w http.ResponseWriter, r *http.Request) {
	// Parse and Validate request
	id := mux.Vars(r)["id"]
	if id=="" {
		utils.WriteError(fmt.Errorf("missing id to find"), w, 400)
		return;
	}
	// Deal with the request in application layer
	doc, err := h.service.GetBaseDocumentByID(r.Context(), id)
	if err != nil || doc == nil {
		fmt.Println(doc)
		utils.WriteError(fmt.Errorf("could not find document: %v", err), w, 500)
		return;
	}

	utils.WriteJson(transformers.ToBaseModelResponse(*doc), w, 200)
}


// handleGetSince handles the request for getting documents after a certain date
func (h servicePort) handleGetSince(w http.ResponseWriter, r *http.Request) {
	// Parse and Validate request
	dateString := mux.Vars(r)["date"]
	if dateString=="" {
		utils.WriteError(fmt.Errorf("missing date to find"), w, 400)
		return;
	}
	date, err := time.Parse("2006-01-02T15:04:05Z", dateString)
	if err != nil {
		utils.WriteError(fmt.Errorf("invalid date sent: %s", dateString), w, 400)
		return;
	}
	// Deal with the request in application layer
	docs, err := h.service.GetAllCreatedSince(r.Context(), date)
	if err != nil || docs == nil {
		fmt.Println(docs)
		utils.WriteError(fmt.Errorf("could not find documents: %v", err), w, 500)
		return;
	}

	utils.WriteJson(transformers.ToBaseModelResponseArray(docs), w, 200)
}
