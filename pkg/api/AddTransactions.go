package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка чтения ответа"))
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ошибка парсинга"))
		return
	}
	ctx := context.Background()
	err = db.AddTransaction(ctx, &transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("оШибка добавления"))
		return
	}
}
