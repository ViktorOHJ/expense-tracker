package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/api"
	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
	"github.com/ViktorOHJ/expense-tracker/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteHandler_Success(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	// Mock с userID
	mockDB.On("DeleteTransaction", mock.Anything, 1, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)

	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var actualResp models.SuccessResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Contains(t, actualResp.Message, "transaction with id: 1 successfully deleted")

	mockDB.AssertExpectations(t)
}

func TestDeleteHandler_Unauthorized(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)

	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	mockDB.AssertExpectations(t)
}

func TestDeleteHandler_NotFound(t *testing.T) {
	mockDB := new(mocks.DB)
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)
	mockDB.On("DeleteTransaction", mock.Anything, 1, 1).Return(db.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)
	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "transaction not found or access denied", actualResp.Message)

	mockDB.AssertExpectations(t)
}

func TestDeleteHandler_BadRequest(t *testing.T) {
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
		{"Empty ID", "/transaction/?id=", http.StatusBadRequest, "id cannot be empty"},
		{"Invalid ID", "/transaction/?id=abc", http.StatusBadRequest, "id must be a positive number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.query, nil)
			claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
			ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
			req = req.WithContext(ctx)
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
	jwtService := auth.NewJWTService("test-secret")
	passwordService := auth.NewPasswordService()
	s := api.NewServer(mockDB, jwtService, passwordService)

	mockDB.On("DeleteTransaction", mock.Anything, 1, 1).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/transaction/?id=1", nil)

	claims := &auth.Claims{UserID: 1, Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), api.UserContextKey, claims)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	s.DeleteHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var actualResp models.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&actualResp)
	assert.NoError(t, err)

	assert.Equal(t, "failed to delete transaction", actualResp.Message)

	mockDB.AssertExpectations(t)
}
