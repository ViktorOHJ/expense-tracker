package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	// Mock the expected behavior
	mockDB.On("DeleteTransaction", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)
	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualResp models.SuccessResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Contains(t, actualResp.Message, "transaction with id: 1 successfully deleted")

	mockDB.AssertExpectations(t)
}

func TestDeleteHandler_NotFound(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	// Mock the expected behavior for not found
	mockDB.On("DeleteTransaction", mock.Anything, 1).Return(db.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)
	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "transaction not found", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestDeleteHandler_BadRequest(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	tests := []struct {
		name       string
		query      string
		statusCode int
		message    string
	}{
		{"Empty ID", "/transaction/?id=", http.StatusBadRequest, "id cannot be empty"},
		{"Invalid ID", "/transaction/?id=abc", http.StatusBadRequest, "id must be a positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.query, nil)
			rr := httptest.NewRecorder()

			s.DeleteHandler(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)

			var actualResp models.ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&actualResp)
			assert.NoError(t, err)

			assert.Equal(t, tt.message, actualResp.Message)
		})
	}
}

func TestDeleteHandler_DBError(t *testing.T) {
	mockDB := new(mocks.DB)
	s := api.NewServer(mockDB)

	// Mock the expected behavior for a database error
	mockDB.On("DeleteTransaction", mock.Anything, 1).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)
	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "failed to delete transaction", actualResp.Message)

	mockDB.AssertExpectations(t)
}
