package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	input := models.Transaction{
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		Note:       "Lunch",
	}

	returned := models.Transaction{
		ID:         1,
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		UserID:     1,
		Note:       "Lunch",
		CreatedAt:  time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
	}

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	mockDB.On("CheckCategory", mock.Anything, 1, input.CategoryID).Return(true, nil)
	mockDB.On("AddTransaction", mock.Anything, 1, mock.MatchedBy(func(tx *models.Transaction) bool {
		return tx.Amount == input.Amount && tx.Note == input.Note
	})).Return(returned, nil)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	claims := &auth.Claims{
		UserID: 1,
		Email:  "test@example.com",
	}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var actualResp models.SuccessResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	dataBytes, err := json.Marshal(actualResp.Data)
	assert.NoError(t, err)

	var actualTransaction models.Transaction
	err = json.Unmarshal(dataBytes, &actualTransaction)
	assert.NoError(t, err)

	assert.Equal(t, returned, actualTransaction)
	assert.Equal(t, "transaction added successfully", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestAddHandler_Unauthorized(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	input := models.Transaction{
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		Note:       "Lunch",
	}

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var actualResp models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "unauthorized", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestAddHandler_InvalidJSON(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	invalidJSON := `{"IsIncome": true, "Amount": 100, "CategoryID": 2, "Note": "Lunch"`

	req := httptest.NewRequest(http.MethodPost, "/transactions?", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	claims := &auth.Claims{
		UserID: 1,
		Email:  "test@example.com",
	}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "invalid json format", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestAddHandler_CategoryNotFound(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	input := models.Transaction{
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		Note:       "Lunch",
	}

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	mockDB.On("CheckCategory", mock.Anything, 1, input.CategoryID).Return(false, nil)

	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var actualResp models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "category does not exist or access denied", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestAddHandler_AmountZeroOrNegative(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	input := models.Transaction{
		IsIncome:   false,
		Amount:     0, // Invalid amount
		CategoryID: 2,
		Note:       "Lunch",
	}

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/transactions?", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var actualResp models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "amount must be greater than 0", actualResp.Message)

	mockDB.AssertExpectations(t)
}
func TestAddHandler_DBError(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	input := models.Transaction{
		IsIncome:   false,
		Amount:     100,
		CategoryID: 2,
		Note:       "Lunch",
	}

	body, err := json.Marshal(input)
	assert.NoError(t, err)

	mockDB.On("CheckCategory", mock.Anything, 1, input.CategoryID).Return(true, nil)
	mockDB.On("AddTransaction", mock.Anything, 1, mock.Anything).Return(models.Transaction{}, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/transactions?", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.AddHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var actualResp models.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "error adding transaction", actualResp.Message)

	mockDB.AssertExpectations(t)
}
