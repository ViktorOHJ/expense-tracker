package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
)

func TestTransactionByIdHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	expectedTransaction := models.Transaction{
		ID:         1,
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		Note:       "Lunch",
		CreatedAt:  time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	mockDB.On("GetTransactionByID", mock.Anything, mock.Anything, 1).Return(expectedTransaction, nil)

	req := httptest.NewRequest(http.MethodGet, "/transaction?id=1", nil)
	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	s.TransactionByIdHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualResp models.SuccessResponse

	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	dataBytes, err := json.Marshal(actualResp.Data)
	assert.NoError(t, err)

	var actualTransaction models.Transaction
	err = json.Unmarshal(dataBytes, &actualTransaction)
	assert.NoError(t, err)

	assert.Equal(t, expectedTransaction, actualTransaction)

	assert.Contains(t, actualResp.Message, "transaction with id: 1 successfully retrieved")

	mockDB.AssertExpectations(t)
}

func TestTransactionByIdHandler_NotFound(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	mockDB.On("GetTransactionByID", mock.Anything, mock.Anything, 1).Return(models.Transaction{}, db.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/transaction?id=1", nil)
	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	s.TransactionByIdHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "transaction with id 1 not found or access denied", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestTransactionByIdHandler_BadRequest(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	tests := []struct {
		name       string
		query      string
		statusCode int
		message    string
	}{
		{"Empty ID", "/transaction?id=", http.StatusBadRequest, "id cannot be empty"},
		{"Invalid ID", "/transaction?id=abc", http.StatusBadRequest, "id must be a positive number"},
		{"Negative ID", "/transaction?id=-1", http.StatusBadRequest, "id must be a positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
			ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			s.TransactionByIdHandler(rr, req)

			assert.Equal(t, tt.statusCode, rr.Code)

			var actualResp models.ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&actualResp)
			assert.NoError(t, err)

			assert.Equal(t, tt.message, actualResp.Message)
		})
	}
}
