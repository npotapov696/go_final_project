package api

import (
	"net/http"

	"go1f/pkg/db"
)

// TasksResp обёртка над слайсом задач для удобства вывода в json-фомате.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

var maxEntries = 10 // максимальное количество выводимых записей

// tasksHandler обрабатывает GET-запрос на возврат списка задач, отсортированных по степени актуальности
// во времени, в json-формате. Количество ограничено значением maxEntries. В случае неудачи возвращает ошибку в
// json-формате.
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(maxEntries, search)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJsonErr(w, err)
		return
	}
	writeJson(w, TasksResp{
		Tasks: tasks,
	})

}
