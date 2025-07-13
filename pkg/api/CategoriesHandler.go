package api

import (
	"encoding/json"
	"io"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func CategoriesHandler(w http.ResponseWriter, r *http.Request) {
	category := models.Category{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "request reading error")
		return
	}
	err = json.Unmarshal(body, &category)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "invalid json format")
		return
	}
	if category.Name == "" {
		JsonError(w, http.StatusBadRequest, "category name cannot be empty")
		return
	}

	ctx := r.Context()
	category, err = db.AddCategory(ctx, &category)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error adding category")
		return
	}
	resp := models.SuccessResponse{
		Message: "category added successfuly",
		Data:    category,
	}
	JsonResponse(w, http.StatusCreated, resp)
}
