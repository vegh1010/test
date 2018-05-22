package merchant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"example.com/test/pkg/db"
	"example.com/test/pkg/env"
	"example.com/test/pkg/logger"
	"example.com/test/pkg/model/merchant"
	"example.com/test/pkg/model/modelinit"
	"example.com/test/pkg/resperror"
	"example.com/test/pkg/txcontext"
)

// environment
var e = env.NewEnv()

// logger
var l = logger.NewLogger(e)

// database
var d = db.NewDB(l, e)

var recIDs []string

// setup test data
func setup(t *testing.T) (merchant.Record, func()) {

	modelinit.PrepareStatements(d)

	// tx
	tx, _ := d.Beginx()

	// merchant
	mm, _ := merchant.NewModel(e, l, tx)
	mr := mm.NewRecord()
	mr.Name = "Test merchant"
	mr.ShortName = "Test merchant"
	mr.DBAName = "Test merchant"
	mr.CountryID = "US"
	mr.TimezoneID = "America/Denver"
	mr.Status = "active"
	mm.Create(&mr)

	tx.Commit()

	// teardown test data
	teardown := func() {

		// tx
		tx, _ = d.Beginx()

		// remove merchant records
		// - created in tests
		mm, _ := merchant.NewModel(e, l, tx)
		for i := range recIDs {
			mm.Remove(recIDs[i])
		}

		// remove merchant
		mm.Remove(mr.ID)

		tx.Commit()
	}

	return mr, teardown
}

// testMerchantPost - tests creating a merchant
func TestMerchantPost(t *testing.T) {

	_, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	// handler
	h := NewHandler(e, l)

	// recorder
	rr := httptest.NewRecorder()

	// data
	er := Request{
		Data: &Data{
			Name:      "Test merchant name",
			ShortName: "Test merchant short name",
			DBAName:   "Test merchant dba name",
			Country:   "US",
			Timezone:  "America/Chicago",
			Status:    "active",
		},
	}

	erj, _ := json.Marshal(er)

	url := h.GetPath()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(erj))
	if err != nil {
		t.Fatal(err)
	}

	// set tx context
	req = txcontext.SetContext(req, tx)

	invokeHandler(http.HandlerFunc(h.Post), h.GetPath(), rr, req)

	// test status
	assert.Equal(t, http.StatusOK, rr.Code, "Create response status code is OK")

	// test body
	res := Response{}
	json.NewDecoder(rr.Body).Decode(&res)

	assert.NotEmpty(t, res.Data.ID, "Merchant ID is not nil")
	assert.Equal(t, "Test merchant name", res.Data.Name, "Merchant Name equals expected")
	assert.Equal(t, "Test merchant short name", res.Data.ShortName, "Merchant ShortName equals expected")
	assert.Equal(t, "Test merchant dba name", res.Data.DBAName, "Merchant DBAName equals expected")
	assert.Equal(t, "US", res.Data.Country, "Merchant CountryID equals expected")
	assert.Equal(t, "America/Chicago", res.Data.Timezone, "Merchant TimezoneID equals expected")
	assert.Equal(t, "inactive", res.Data.Status, "Merchant Status equals expected")
	assert.NotEmpty(t, res.Data.CreatedAt, "Merchant CreatedAt is not nil")

	recIDs = append(recIDs, res.Data.ID)

	tx.Rollback()
}

// testMerchantGet -
func TestMerchantGet(t *testing.T) {

	mr, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	// merchant handler
	h := NewHandler(e, l)

	// recorder
	rr := httptest.NewRecorder()

	// request
	url := h.GetPath() + "/" + mr.ID

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// set tx context
	req = txcontext.SetContext(req, tx)

	invokeHandler(http.HandlerFunc(h.Get), h.GetPath()+"/{id}", rr, req)

	// test status
	assert.Equal(t, http.StatusOK, rr.Code, "Get response status code is OK")

	// test body
	res := Response{}
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, mr.Name, res.Data.Name, "Merchant Name equals expected")
	assert.Equal(t, mr.ShortName, res.Data.ShortName, "Merchant ShortName equals expected")
	assert.Equal(t, mr.DBAName, res.Data.DBAName, "Merchant DBAName equals expected")
	assert.Equal(t, mr.CountryID, res.Data.Country, "Merchant CountryID equals expected")
	assert.Equal(t, mr.TimezoneID, res.Data.Timezone, "Merchant TimezoneID equals expected")
	assert.Equal(t, mr.Status, res.Data.Status, "Merchant Status equals expected")
	assert.NotEmpty(t, res.Data.CreatedAt, "Merchant CreatedAt is not nil")
	assert.Empty(t, res.Data.UpdatedAt, "Merchant UpdatedAt is empty")

	tx.Rollback()
}

func TestMerchantUpdate(t *testing.T) {

	mr, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	// merchant handler
	h := NewHandler(e, l)

	// recorder
	rr := httptest.NewRecorder()

	// data
	er := Request{
		Data: &Data{
			Name:      "Test update merchant name",
			ShortName: "Test update merchant short name",
			DBAName:   "Test update merchant dba name",
			Country:   "AU",
			Timezone:  "Australia/Sydney",
			Status:    "inactive",
		},
	}

	erj, _ := json.Marshal(er)

	url := h.GetPath() + "/" + mr.ID

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(erj))
	if err != nil {
		t.Fatal(err)
	}

	// set tx context
	req = txcontext.SetContext(req, tx)

	invokeHandler(http.HandlerFunc(h.Put), h.GetPath()+"/{id}", rr, req)

	// test status
	assert.Equal(t, http.StatusOK, rr.Code, "Update response status code is OK")

	// test body
	res := Response{}
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, er.Data.Name, res.Data.Name, "Merchant Name equals expected")
	assert.Equal(t, er.Data.ShortName, res.Data.ShortName, "Merchant ShortName equals expected")
	assert.Equal(t, er.Data.DBAName, res.Data.DBAName, "Merchant DBAName equals expected")
	assert.Equal(t, er.Data.Country, res.Data.Country, "Merchant CountryID equals expected")
	assert.Equal(t, er.Data.Timezone, res.Data.Timezone, "Merchant TimezoneID equals expected")
	assert.Equal(t, er.Data.Status, res.Data.Status, "Merchant status equals expected")
	assert.NotEmpty(t, res.Data.CreatedAt, "Merchant CreatedAt is not nil")
	assert.NotEmpty(t, res.Data.UpdatedAt, "Merchant UpdatedAt is not empty")

	tx.Rollback()
}

func testMerchantDelete(t *testing.T, recID string) {

	_, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	// merchant handler
	h := NewHandler(e, l)

	// recorder
	rr := httptest.NewRecorder()

	url := h.GetPath() + "/" + recID

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// set tx context
	req = txcontext.SetContext(req, tx)

	invokeHandler(http.HandlerFunc(h.Delete), h.GetPath()+"/{id}", rr, req)

	// test status
	assert.Equal(t, http.StatusOK, rr.Code, "Delete response status code is OK")

	tx.Rollback()
}

func testMerchantPostValidation(t *testing.T) {

	_, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	h := NewHandler(e, l)

	rr := httptest.NewRecorder()

	tx, _ = d.Beginx()

	// data
	er := Request{
		Data: &Data{},
	}

	erj, _ := json.Marshal(er)

	url := h.GetPath()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(erj))
	if err != nil {
		t.Fatal(err)
	}

	req = txcontext.SetContext(req, tx)

	h.Post(rr, req)

	// test status
	assert.Equal(t, rr.Code, http.StatusBadRequest, "Create response status code is Bad Request")

	// test error
	res := resperror.Response{}
	json.NewDecoder(rr.Body).Decode(&res)

	assert.NotNil(t, res.Error, "Validation errors are not nil")

	assert.NotNil(t, res.Error.Code, "Error Code is not nil")
	assert.NotNil(t, res.Error.Title, "Error Title is not nil")
	assert.NotNil(t, res.Error.Detail, "Error Detail is not nil")

	tx.Rollback()
}

func testMerchantGetCollection(t *testing.T) {

	_, teardown := setup(t)
	defer teardown()

	// tx
	tx, _ := d.Beginx()

	h := NewHandler(e, l)

	rr := httptest.NewRecorder()

	// request
	url := h.GetPath()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// set tx context
	req = txcontext.SetContext(req, tx)

	h.GetCollection(rr, req)

	// test status
	assert.Equal(t, http.StatusOK, rr.Code, "Get collection response status is OK")

	tx.Rollback()
}

func testMerchantPostFailure(t *testing.T) {

	_, teardown := setup(t)
	defer teardown()

	h := NewHandler(e, l)

	type TestCase struct {
		Req        Request
		StatusCode int
		ErrorCode  int
	}

	for testNumber, tc := range []TestCase{
		TestCase{
			Req: Request{
				Data: &Data{
					Name:      "Test merchant name",
					ShortName: "Test merchant short name",
					DBAName:   "Test merchant dba name",
					Country:   "ZZ",
					Timezone:  "America/Chicago",
					Status:    "active",
				},
			},
			StatusCode: http.StatusBadRequest,
			ErrorCode:  resperror.ErrorInvalidCountry.Code,
		},
		TestCase{
			Req: Request{
				Data: &Data{
					Name:      "Test merchant name",
					ShortName: "Test merchant short name",
					DBAName:   "Test merchant dba name",
					Country:   "US",
					Timezone:  "Moon/OtherSide",
					Status:    "active",
				},
			},
			StatusCode: http.StatusBadRequest,
			ErrorCode:  resperror.ErrorInvalidTimezone.Code,
		},
	} {

		rr := httptest.NewRecorder()

		erj, _ := json.Marshal(tc.Req)

		tx, _ := d.Beginx()

		// url
		url := "/api/merchants"

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(erj))

		// set tx context
		req = txcontext.SetContext(req, tx)

		if err != nil {
			t.Fatal(err)
		}

		h.Post(rr, req)

		// test status
		assert.Equal(t, tc.StatusCode, rr.Code, fmt.Sprintf("Create merchant test : %d status code is: %d as expected got %v - %v", testNumber, tc.StatusCode, rr.Code, rr.Body))

		res := resperror.Response{}
		json.NewDecoder(rr.Body).Decode(&res)

		assert.Equal(t, tc.ErrorCode, res.Error.Code, "Expected error code %d, got %d", tc.ErrorCode, res.Error.Code)

		tx.Rollback()
	}
}

func testMerchantPutFailure(t *testing.T) {

	mr, teardown := setup(t)
	defer teardown()

	h := NewHandler(e, l)

	type TestCase struct {
		Req        Request
		StatusCode int
		ErrorCode  int
	}

	for testNumber, tc := range []TestCase{
		// InvalidCountry
		TestCase{
			Req: Request{
				Data: &Data{
					Name:      "Test merchant name",
					ShortName: "Test merchant short name",
					DBAName:   "Test merchant dba name",
					Country:   "ZZ",
					Timezone:  "America/Chicago",
					Status:    "active",
				},
			},
			StatusCode: http.StatusBadRequest,
			ErrorCode:  resperror.ErrorInvalidCountry.Code,
		},
		// InvalidTimezone
		TestCase{
			Req: Request{
				Data: &Data{
					Name:      "Test merchant name",
					ShortName: "Test merchant short name",
					DBAName:   "Test merchant dba name",
					Country:   "US",
					Timezone:  "Moon/OtherSide",
					Status:    "active",
				},
			},
			StatusCode: http.StatusBadRequest,
			ErrorCode:  resperror.ErrorInvalidTimezone.Code,
		},
	} {

		rr := httptest.NewRecorder()

		erj, _ := json.Marshal(tc.Req)

		tx, _ := d.Beginx()

		// url
		url := "/api/merchants/" + mr.ID

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(erj))

		// set tx context
		req = txcontext.SetContext(req, tx)

		if err != nil {
			t.Fatal(err)
		}

		invokeHandler(http.HandlerFunc(h.Put), h.GetPath()+"/{id}", rr, req)

		// test status
		assert.Equal(t, tc.StatusCode, rr.Code, fmt.Sprintf("Update merchant test : %d status code is: %d as expected got %v - %v", testNumber, tc.StatusCode, rr.Code, rr.Body))

		res := resperror.Response{}
		json.NewDecoder(rr.Body).Decode(&res)

		assert.Equal(t, tc.ErrorCode, res.Error.Code, "Expected error code %d, got %d", tc.ErrorCode, res.Error.Code)

		tx.Rollback()
	}
}

// invokeHandler is used to test parametized requests
func invokeHandler(handler http.Handler, path string, w http.ResponseWriter, r *http.Request) {

	router := mux.NewRouter()

	router.Path(path).Handler(handler)

	router.ServeHTTP(w, r)
}
