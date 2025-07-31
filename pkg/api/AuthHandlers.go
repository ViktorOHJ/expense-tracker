package api

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "request reading error")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		JsonError(w, http.StatusBadRequest, "invalid json format")
		return
	}

	// Валидация
	if !isValidEmail(req.Email) {
		JsonError(w, http.StatusBadRequest, "invalid email format")
		return
	}

	if len(req.Password) < 6 {
		JsonError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	// Проверяем, существует ли пользователь
	_, err = s.db.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		JsonError(w, http.StatusConflict, "user already exists")
		return
	}

	// Хешируем пароль
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error processing password")
		return
	}

	// Создаем пользователя
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	createdUser, err := s.db.CreateUser(r.Context(), &user)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	// Генерируем токен
	token, err := s.jwtService.GenerateToken(createdUser.ID, createdUser.Email)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  createdUser,
	}

	JsonResponse(w, http.StatusCreated, models.SuccessResponse{
		Message: "user registered successfully",
		Data:    response,
	})
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "request reading error")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		JsonError(w, http.StatusBadRequest, "invalid json format")
		return
	}

	// Получаем пользователя
	user, err := s.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == db.ErrNotFound {
			JsonError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		JsonError(w, http.StatusInternalServerError, "error retrieving user")
		return
	}

	// Проверяем пароль
	if !s.passwordService.CheckPassword(user.Password, req.Password) {
		JsonError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Генерируем токен
	token, err := s.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error generating token")
		return
	}

	response := models.AuthResponse{
		Token: token,
		User:  user,
	}

	JsonResponse(w, http.StatusOK, models.SuccessResponse{
		Message: "login successful",
		Data:    response,
	})
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}
