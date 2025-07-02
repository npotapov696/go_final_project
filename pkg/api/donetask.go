package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

// doneHandler обрабатывает POST-запрос по переданному в URL "id" на изменение даты задачи на
// актуальную в базе данных, либо на удаление, если правило задачи отсутствует. В случае успешного
// выполнения возвращает пустой json. В случае неудачи возвращает ошибку в json-формате.
func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	if len(task.Repeat) == 0 {
		err := db.DeleteTask(task.ID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJsonErr(w, err)
			return
		}
		writeJson(w, map[string]interface{}{})
		return
	}
	task.Date, err = nextDate(time.Now().UTC(), task.Date, task.Repeat)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
	}
	err = db.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	writeJson(w, map[string]interface{}{})
}
