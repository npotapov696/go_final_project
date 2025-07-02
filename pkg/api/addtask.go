package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
)

// addTaskHandler обрабатывает POST-запрос, в теле которого передан экземпляр структуры задачи в
// json-формате, на добавление этой задачи в таблицу базы данных. В случае успеха возвращает "id"
// в json-формате, в случае неудачи - ошибку в json-формате.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer
	var jsID JsonID
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
	id, err := db.AddTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	jsID.ID = strconv.Itoa(int(id))
	w.WriteHeader(http.StatusOK)
	writeJson(w, jsID)
}

// checkDate проверяет на корректность даты задачи, переданной в task.
func checkDate(task *db.Task) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if len(task.Date) == 0 {
		task.Date = now.Format(db.DateString)
	}
	t, err := time.Parse(db.DateString, task.Date)
	if err != nil {
		return err
	}
	if !t.After(now) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(db.DateString)
			return nil
		}
		task.Date, err = nextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	return nil
}
