package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go1f/pkg/db"
)

// updateTaskHandler обрабатывает PUT-запрос, в теле которого передан экземпляр структуры задачи в
// json-формате, на обновление полей таблицы базы данных, соответствующих "id".
// В случае успеха возвращает пустой json, в случае неудачи - ошибку в json-формате.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, fmt.Errorf("не указан заголовок задачи"))
		return
	}
	if err = checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, map[string]interface{}{})
}
