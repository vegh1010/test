package merchant

import (
	"net/http"

	"github.com/rs/zerolog"

	"example.com/test/pkg/env"
	"example.com/test/pkg/handler"
	"example.com/test/pkg/model/merchant"
	"example.com/test/pkg/resperror"
)

// Data -
type Data struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	DBAName   string `json:"dba_name"`
	Country   string `json:"country"`
	Timezone  string `json:"timezone"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Response -
type Response struct {
	Data *Data `json:"data"`
}

// CollectionResponse -
type CollectionResponse struct {
	Data []*Data `json:"data"`
}

// Request -
type Request struct {
	Data *Data `json:"data"`
}

// Handler -
type Handler struct {
	handler.Base
}

// NewHandler -
func NewHandler(e *env.Env, l zerolog.Logger) handler.Handler {
	h := Handler{
		handler.Base{
			Path:            "/api/merchants",
			Unauthenticated: false, // Requires authentication
			Unauthorized:    false, // Requires authorization
			Versioned:       true,
			Env:             e,
			Logger:          l,
			LockResources: map[string]map[string]string{
				http.MethodPut: {"merchant": "id"},
			},
		},
	}
	return &h
}

// Get -
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {

	// logger
	log := h.Logger

	ms, params, err := h.PreHandlerChecks(r)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	// model
	m, err := ms.GetMerchantModel()
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Get with params %v", params)

	// get
	recs, err := m.GetByParam(params)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	if len(recs) != 1 {
		// not found
		h.SendErrorResponse(w, r, resperror.ErrorNotFound)
		return
	}

	// record
	rec := recs[0]

	res := Response{
		Data: &Data{
			ID:        rec.ID,
			Name:      rec.Name,
			ShortName: rec.ShortName,
			DBAName:   rec.DBAName,
			Country:   rec.CountryID,
			Timezone:  rec.TimezoneID,
			Status:    rec.Status,
			CreatedAt: rec.CreatedAt,
			UpdatedAt: rec.UpdatedAt.String,
		},
	}

	h.DebugStruct("Get Response", res)

	h.SendResponse(w, r, &res)

	log.Debug().Msgf("Merchant fetched OK")
}

// GetCollection -
func (h *Handler) GetCollection(w http.ResponseWriter, r *http.Request) {

	// logger
	log := h.Logger

	ms, params, err := h.PreHandlerChecks(r)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	// model
	m, err := ms.GetMerchantModel()
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("GetCollection with params %v", params)

	recs, _ := m.GetByParam(params)
	if recs != nil {

		var ed []*Data

		for _, rec := range recs {
			ed = append(ed, &Data{
				ID:        rec.ID,
				Name:      rec.Name,
				ShortName: rec.ShortName,
				DBAName:   rec.DBAName,
				Country:   rec.CountryID,
				Timezone:  rec.TimezoneID,
				Status:    rec.Status,
				CreatedAt: rec.CreatedAt,
				UpdatedAt: rec.UpdatedAt.String,
			})
		}

		res := CollectionResponse{
			Data: ed,
		}

		h.DebugStruct("Get Response", res)

		h.SendResponse(w, r, &res)
	}

	log.Debug().Msgf("Merchant fetched OK")
}

// Post -
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {

	// logger
	log := h.Logger

	ms, params, err := h.PreHandlerChecks(r)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Post with params %v", params)

	// decode request body
	req := Request{}
	err = h.DecodeRequest(r, &req)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Post with data %v", req)

	// validate
	verr := req.Validate()
	if verr != nil {
		h.SendErrorResponse(w, r, verr)
		return
	}

	// model
	m, err := ms.GetMerchantModel()
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	// record
	// example: rec.XxxxID = params["xxx_id"].(string)
	// - NOTE
	// - This is where we would decide which
	// - properties can and cannot be set
	// - based on role?
	rec := m.NewRecord()
	rec.Name = req.Data.Name
	rec.ShortName = req.Data.ShortName
	rec.DBAName = req.Data.DBAName
	rec.CountryID = req.Data.Country
	rec.TimezoneID = req.Data.Timezone

	log.Debug().Msgf("Validate with record %v", rec)
	vrec, err := m.ValidateRecord(&rec)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	if vrec.CountryID.Bool == false {
		h.SendErrorResponse(w, r, resperror.ErrorInvalidCountry)
		return
	}
	if vrec.TimezoneID.Bool == false {
		h.SendErrorResponse(w, r, resperror.ErrorInvalidTimezone)
		return
	}

	log.Debug().Msgf("Create with record %v", rec)

	// create
	err = m.Create(&rec)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	if rec.ID != "" {
		res := Response{
			Data: &Data{
				ID:        rec.ID,
				Name:      rec.Name,
				ShortName: rec.ShortName,
				DBAName:   rec.DBAName,
				Country:   rec.CountryID,
				Timezone:  rec.TimezoneID,
				Status:    rec.Status,
				CreatedAt: rec.CreatedAt,
				UpdatedAt: rec.UpdatedAt.String,
			},
		}

		h.DebugStruct("Post Response", res)

		h.SendResponse(w, r, &res)
	}

	log.Debug().Msgf("Merchant created OK")
}

// Put -
func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {

	// logger
	log := h.Logger

	ms, params, err := h.PreHandlerChecks(r)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	// model
	m, err := ms.GetMerchantModel()
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Put with params %v", params)

	// decode request body
	req := Request{}
	err = h.DecodeRequest(r, &req)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Put with data %v", req)

	// validate
	verr := req.Validate()
	if verr != nil {
		h.SendErrorResponse(w, r, verr)
		return
	}

	// get current record
	recs, _ := m.GetByParam(params)
	if len(recs) != 1 || recs[0].ID != params["id"].(string) {
		// not found
		h.SendErrorResponse(w, r, resperror.ErrorNotFound)
		return
	}

	// record
	rec := recs[0]

	if rec.Status == merchant.StatusTerminated {
		h.SendErrorResponse(w, r, resperror.ErrTerminatedMerchantCannotBeModified)
		return
	}

	// update record properties
	// - NOTE
	// - This is where we would decide which
	// - properties can and cannot be updated
	// - based on role?
	rec.Name = req.Data.Name
	rec.ShortName = req.Data.ShortName
	rec.DBAName = req.Data.DBAName
	rec.CountryID = req.Data.Country
	rec.TimezoneID = req.Data.Timezone
	rec.Status = req.Data.Status

	log.Debug().Msgf("Validate with record %v", rec)
	vrec, err := m.ValidateRecord(rec)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	if vrec.CountryID.Bool == false {
		h.SendErrorResponse(w, r, resperror.ErrorInvalidCountry)
		return
	}
	if vrec.TimezoneID.Bool == false {
		h.SendErrorResponse(w, r, resperror.ErrorInvalidTimezone)
		return
	}

	// update
	err = m.Update(rec)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	if rec.ID != "" {
		res := Response{
			Data: &Data{
				ID:        rec.ID,
				Name:      rec.Name,
				ShortName: rec.ShortName,
				DBAName:   rec.DBAName,
				Country:   rec.CountryID,
				Timezone:  rec.TimezoneID,
				Status:    rec.Status,
				CreatedAt: rec.CreatedAt,
				UpdatedAt: rec.UpdatedAt.String,
			},
		}

		h.DebugStruct("Put Response", res)

		h.SendResponse(w, r, &res)
	}

	log.Debug().Msgf("Merchant updated OK")
}

// Delete -
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	// logger
	log := h.Logger

	ms, params, err := h.PreHandlerChecks(r)
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Delete merchant with params %v", params)

	// model
	m, err := ms.GetMerchantModel()
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	// get current record
	recs, _ := m.GetByParam(params)
	if len(recs) != 1 || recs[0].ID != params["id"].(string) {
		// not found
		h.SendErrorResponse(w, r, resperror.ErrorNotFound)
		return
	}

	// delete
	err = m.Delete(params["id"].(string))
	if err != nil {
		h.SendErrorResponse(w, r, err)
		return
	}

	log.Debug().Msgf("Merchant deleted OK")

	h.SendResponse(w, r, nil)
}
