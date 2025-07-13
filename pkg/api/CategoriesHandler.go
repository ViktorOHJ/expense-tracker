package api

import (
	"context"
	"encoding/json"
	"fmt"
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
	ctx := context.Background()
	id, err := db.AddCategory(ctx, &category)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error adding category")
		return
	}
	category.ID = id
	resp := models.SuccessResponse{
		Message: fmt.Sprintf("category with id: %d added successfuly", id),
		Data:    category,
	}
	JsonResponse(w, http.StatusCreated, resp)
}
