package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/test/pkg/txcontext"

	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"database/sql"
	"strings"

	"example.com/test/pkg/env"
	"example.com/test/pkg/modelstore"
	"example.com/test/pkg/resperror"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"gopkg.in/olivere/elastic.v6"
)

// Handler -
type Handler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Post(w http.ResponseWriter, r *http.Request)
	Put(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetCollection(w http.ResponseWriter, r *http.Request)
	PutCollection(w http.ResponseWriter, r *http.Request)
	DeleteCollection(w http.ResponseWriter, r *http.Request)
	GetPath() string
	GetUnauthenticated() bool
	GetUnauthorized() bool
	GetVersioned() bool
	GetLogger() zerolog.Logger
	GetLockResources() map[string]map[string]string
}

// LockResource -
type LockResource struct {
	Method    string
	Resources map[string]string
}

// Base -
type Base struct {
	Path            string
	Unauthenticated bool
	Unauthorized    bool
	Versioned       bool
	Env             *env.Env
	Logger          zerolog.Logger
	LockResources   map[string]map[string]string
}

// Params -
type Params map[string]interface{}

// NewModelStore returns a new model store
func (h *Base) NewModelStore(tx *sqlx.Tx) (*modelstore.ModelStore, error) {

	ms, err := modelstore.NewModelStore(h.Env, h.Logger, tx)
	if err != nil {
		log.Error().Msgf("Error getting new modelstore: %v", err)
		return nil, err
	}

	return ms, nil
}

// Get -
func (h *Base) Get(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// Post -
func (h *Base) Post(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// Put -
func (h *Base) Put(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// Delete -
func (h *Base) Delete(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// GetCollection -
func (h *Base) GetCollection(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// PutCollection -
func (h Base) PutCollection(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// DeleteCollection -
func (h *Base) DeleteCollection(w http.ResponseWriter, r *http.Request) {

	// log
	log := h.Logger

	log.Error().Msgf("Method not implemented for path %s", r.RequestURI)
	http.Error(w, http.StatusText(http.StatusNotFound),
		http.StatusNotFound)
}

// DebugStruct -
func (h *Base) DebugStruct(msg string, rec interface{}) {

	// logger
	log := h.Logger

	log.Debug().Msgf(msg + " " + spew.Sdump(rec))
}

// ValidateParams -
func (h *Base) ValidateParams(vars map[string]string) (params Params, rerr *resperror.Data) {

	// logger
	log := h.Logger

	// validate vars
	params = Params{}

	for k, v := range vars {
		log.Debug().Msgf("Validate param k [%s] v [%s]", k, v)

		if v == "" || v == "{"+k+"}" {
			log.Debug().Msgf("Validate error for param %s", k)

			return nil, resperror.ValidationRequired(k)
		}

		params[k] = v
	}

	return params, nil
}

// GetPath -
func (h *Base) GetPath() string {
	return h.Path
}

// GetUnauthenticated -
func (h *Base) GetUnauthenticated() bool {
	return h.Unauthenticated
}

// GetUnauthorized -
func (h *Base) GetUnauthorized() bool {
	return h.Unauthorized
}

// GetLockResources -
func (h *Base) GetLockResources() map[string]map[string]string {
	return h.LockResources
}

// GetVersioned -
func (h *Base) GetVersioned() bool {
	return h.Versioned
}

// GetLogger -
func (h *Base) GetLogger() zerolog.Logger {
	return h.Logger
}

// DecodeRequest -
func (h *Base) DecodeRequest(r *http.Request, s interface{}) error {
	if r.Body == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(s)
}

// SendErrorResponse sends an error response to the user.
//
// It calls rollback on any db tx available in the request's context
// and then checks what type of error to respond to the user. If the error
// is an unknown type, a system error is returned.
func (h *Base) SendErrorResponse(w http.ResponseWriter, r *http.Request, e error) error {
	// Rollback the active tx.
	err := h.rollbackTx(r)
	if err != nil {
		h.Logger.Error().Msgf("Failed to rollback tx: %s", err.Error())
	}

	var rerr resperror.Response

	// Default the status code to an internal server error.
	httpcode := http.StatusInternalServerError

	switch e {
	case sql.ErrNoRows:
		rerr.Error = resperror.ErrorNotFound
		return h.sendErrorResponse(w, r, &rerr, http.StatusNotFound)
	}

	// Check what type of error is being returned.
	switch et := e.(type) {
	case *resperror.Data:
		if resperror.IsValidationErr(et.Code) {
			httpcode = http.StatusBadRequest
		}
		if et.Code == resperror.ErrCodeNotFound {
			httpcode = http.StatusNotFound
		}
		rerr.Error = et
	case *json.SyntaxError:
		rerr.Error = resperror.ValidationJSONSyntax(et.Offset)
		httpcode = http.StatusBadRequest
	case *pq.Error:
		// NOTE: https://godoc.org/github.com/lib/pq#Error
		//
		// May be able to extract information based on pq.Error values
		// stored in et var here. Example: et.Code or et.Constraint
		//
		// Benji NOTE: I may not have previously used this purely because
		// of a lib/pq version or actual postgres version incompatibility,
		// however this may be possible now.

		// Current known errors.
		if strings.Contains(e.Error(), "invalid input syntax for uuid") {
			h.Logger.Warn().Msgf("Database error, malformed UUID %v", e)
			// Return a not found error for malformed UUID
			// - Usually means a resource wouldn't be able to be found
			rerr.Error = resperror.ErrorNotFound
			httpcode = http.StatusNotFound
		} else {
			// Return an internal error, bad status
			// if we can't work out what it is.
			h.Logger.Error().Msgf("Database error, %v", e)
			rerr.Error = resperror.SystemErr("An internal error has occurred")
			httpcode = http.StatusBadRequest
		}

	case *elastic.Error:
		if et.Status == http.StatusBadRequest {
			httpcode = http.StatusBadRequest
			rerr.Error = resperror.ValidationErr(fmt.Sprintf("Elasticsearch error: %v", et))
		} else {
			rerr.Error = resperror.SystemErr(fmt.Sprintf("Elasticsearch error: %v", et))
		}
	default:
		// Default to a system error.
		rerr.Error = resperror.SystemErr(fmt.Sprintf("Error: %v", e))
	}

	if rerr.Error.Code == resperror.ErrCodeSystem {
		h.Logger.Error().Msgf("%s", rerr.Error.Error())
	}

	// Send the error response to the user.
	return h.sendErrorResponse(w, r, &rerr, httpcode)
}

// SendErrorResponseWithStatusOK is almost identical to SendErrorResponse, but always returns a status code of 200.
func (h *Base) SendErrorResponseWithStatusOK(w http.ResponseWriter, r *http.Request, e error) error {
	// Rollback the active tx.
	err := h.rollbackTx(r)
	if err != nil {
		h.Logger.Error().Msgf("Failed to rollback tx: %s", err.Error())
	}

	var rerr resperror.Response

	// Default the status code to status ok.
	httpcode := http.StatusOK

	switch e {
	case sql.ErrNoRows:
		rerr.Error = resperror.ErrorNotFound
		return h.sendErrorResponse(w, r, &rerr, httpcode)
	}

	// Check what type of error is being returned.
	switch et := e.(type) {
	case *resperror.Data:
		rerr.Error = et
	case *json.SyntaxError:
		rerr.Error = resperror.ValidationJSONSyntax(et.Offset)
	case *pq.Error:

		// Current known errors.
		if strings.Contains(e.Error(), "invalid input syntax for uuid") {
			h.Logger.Warn().Msgf("Database error, malformed UUID %v", e)
			// Return a not found error for malformed UUID
			// - Usually means a resource wouldn't be able to be found
			rerr.Error = resperror.ErrorNotFound
		} else {
			// Return an internal error, bad status
			// if we can't work out what it is.
			h.Logger.Error().Msgf("Database error, %v", e)
			rerr.Error = resperror.SystemErr("An internal error has occurred")
			httpcode = http.StatusInternalServerError
		}

	case *elastic.Error:
		if et.Status == http.StatusBadRequest {
			rerr.Error = resperror.ValidationErr(fmt.Sprintf("Elasticsearch error: %v", et))
		} else {
			rerr.Error = resperror.SystemErr(fmt.Sprintf("Elasticsearch error: %v", et))
			httpcode = http.StatusInternalServerError
		}
	default:
		// Default to a system error.
		rerr.Error = resperror.SystemErr(fmt.Sprintf("Error: %v", e))
		httpcode = http.StatusInternalServerError
	}

	if rerr.Error.Code == resperror.ErrCodeSystem {
		h.Logger.Error().Msgf("%s", rerr.Error.Error())
	}

	// Send the error response to the user.
	return h.sendErrorResponse(w, r, &rerr, httpcode)
}

func (h *Base) sendErrorResponse(w http.ResponseWriter, r *http.Request, rerr *resperror.Response, code int) error {

	// log
	log := h.Logger

	log.Info().Msgf("Sending error response %v", rerr)
	log.Info().Msgf("Sending error response code %d", code)

	// content type json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Status
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(rerr)
}

// SendResponse -
func (h *Base) SendResponse(w http.ResponseWriter, r *http.Request, s interface{}) error {

	// commit tx
	err := h.commitTx(r)
	if err != nil {
		h.Logger.Error().Msgf("Sending error response %v", err)

		res := &resperror.Response{
			Error: resperror.SystemErr("Internal application error"),
		}

		return h.sendErrorResponse(w, r, res, http.StatusInternalServerError)
	}

	// content type json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Status Ok
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(s)
}

func (h *Base) commitTx(r *http.Request) error {

	tx, err := txcontext.GetContext(r)
	if err != nil {
		return fmt.Errorf("Could not commit transaction for context error %v", err)
	}

	return tx.Commit()
}

func (h *Base) rollbackTx(r *http.Request) error {

	tx, err := txcontext.GetContext(r)
	if err != nil {
		return fmt.Errorf("Could not rollback transaction for context error %v", err)
	}

	return tx.Rollback()

}

// PreHandlerChecks -
func (h *Base) PreHandlerChecks(r *http.Request) (*modelstore.ModelStore, Params, error) {

	// get tx from context
	tx, err := txcontext.GetContext(r)
	if err != nil {
		return nil, nil, err
	}

	// modelstore
	ms, err := h.NewModelStore(tx)
	if err != nil {
		return nil, nil, err
	}

	// validate params
	vars := mux.Vars(r)
	params, errs := h.ValidateParams(vars)
	if errs != nil {
		return nil, nil, errs
	}

	return ms, params, nil
}
