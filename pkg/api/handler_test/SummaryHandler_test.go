package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"

	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSummaryHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC)

	expectedSummary := models.Summary{
		TotalIncome:  5000,
		TotalExpense: 3000,
		Balance:      2000,
	}

	mockDB.On("GetSummary", mock.Anything, from, to).Return(expectedSummary, nil)

	req := httptest.NewRequest(http.MethodGet, "/summary?from=2024-01-01&to=2024-12-31", nil)
	rr := httptest.NewRecorder()

	s.SummaryHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualResp models.SuccessResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	dataBytes, err := json.Marshal(actualResp.Data)
	assert.NoError(t, err)

	var actualSummary models.Summary
	err = json.Unmarshal(dataBytes, &actualSummary)
	assert.NoError(t, err)

	assert.Equal(t, expectedSummary.TotalIncome, actualSummary.TotalIncome)
	assert.Equal(t, expectedSummary.TotalExpense, actualSummary.TotalExpense)
	assert.Equal(t, expectedSummary.Balance, actualSummary.Balance)

	mockDB.AssertExpectations(t)
}

func TestSummaryHandler_InvalidDateRange(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	req := httptest.NewRequest(http.MethodGet, "/summary?from=2024-12-31&to=2024-01-01", nil)
	rr := httptest.NewRecorder()

	s.SummaryHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "to date must be after from date", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestSummaryHandler_DBError(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	mockDB.On("GetSummary", mock.Anything, from, to).Return(models.Summary{}, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/summary?from=2024-01-01&to=2024-12-31", nil)
	rr := httptest.NewRecorder()

	s.SummaryHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "error retrieving summary", actualResp.Message)

	mockDB.AssertExpectations(t)
}
