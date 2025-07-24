package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?type=true&category_id=1&from=2024-01-01&to=2024-12-31&limit=10&offset=1", nil)
	rr := httptest.NewRecorder()

	mockDB.On("GetTransactions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, 10, 0).
		Return([]*models.Transaction{
			{ID: 1, Amount: 100, CategoryID: 1},
		}, nil)

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp models.SuccessResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "transactions listed successfully", resp.Message)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)

	var tx []models.Transaction
	err = json.Unmarshal(dataBytes, &tx)
	assert.NoError(t, err)
	assert.Len(t, tx, 1)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_InvalidType(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?type=invalid", nil)
	rr := httptest.NewRecorder()

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid type parameter", resp.Message)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_InvalidDateFormat(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?from=invalid-date", nil)
	rr := httptest.NewRecorder()

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid date format for 'from'", resp.Message)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_InvalidLimit(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?limit=invalid", nil)
	rr := httptest.NewRecorder()

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid limit parameter", resp.Message)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_NoTransactions(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?limit=10&offset=1", nil)
	rr := httptest.NewRecorder()

	mockDB.On("GetTransactions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, 10, 0).
		Return([]*models.Transaction{}, nil)

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp models.SuccessResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "no transactions found", resp.Message)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	var tx []models.Transaction
	err = json.Unmarshal(dataBytes, &tx)
	assert.NoError(t, err)
	assert.Equal(t, len(tx), 0)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_InvalidOffset(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?limit=10&offset=-1", nil)
	rr := httptest.NewRecorder()

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid offset parameter", resp.Message)

	mockDB.AssertExpectations(t)
}

func TestGetHandler_DBError(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/transactions?limit=10&offset=1", nil)
	rr := httptest.NewRecorder()

	mockDB.On("GetTransactions", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, 10, 0).
		Return(nil, assert.AnError)

	s.GetHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
