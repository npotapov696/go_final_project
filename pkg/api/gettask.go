package api

import (
	"net/http"

	"go1f/pkg/db"
)

// getTaskHandler обрабатывает GET-запрос по переданному в URL "id" на возврат задачи в json-формате.
// В случае неудачи возвращает ошибку в json-формате.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, task)
}
