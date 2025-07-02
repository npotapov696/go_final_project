package api

import (
	"net/http"

	"go1f/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE-запрос по переданному в URL "id" на удаление задачи из
// базы данных. В случае успешного выполнения возвращает пустой json. В случае неудачи возвращает
// ошибку в json-формате.
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	writeJson(w, map[string]interface{}{})
}
