// Пакет api реализует маршрутизацию запросов при запуске сервера, а так же описывает хендлер-функции
// для реализации этих запросов, вспомогательные функции, структуры и прочие.
package api

import (
	"encoding/json"
	"net/http"
)

// WebDir содержит путь к файлам фронтэнда.
var WebDir = "./web"

// JsonID обёртка над id для удобства вывода в формате json.
type JsonID struct {
	ID string `json:"id,omitempty"`
}

// JsonErr обёртка над err для удобства вывода в формате json.
type JsonErr struct {
	Err string `json:"error,omitempty"`
}

// Init инициализирует хендлеры.
func Init() {
	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	http.HandleFunc("/api/nextdate", nextDayHandler)

	http.HandleFunc("/api/task", auth(taskHandler))

	http.HandleFunc("/api/tasks", auth(tasksHandler))

	http.HandleFunc("/api/task/done", auth(doneHandler))

	http.HandleFunc("/api/signin", passCheckHandler)

}

// TaskHandler распределяет обращение по адресу в соответствии с методом запроса.
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		getTaskHandler(w, r)
	case r.Method == http.MethodPost:
		addTaskHandler(w, r)
	case r.Method == http.MethodPut:
		updateTaskHandler(w, r)
	case r.Method == http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "метод не поддерживается", http.StatusMethodNotAllowed)
	}

}

// writeJson конвертирует переданные в answer данные в json-формат и записывает в ответ w.
// Если конвертация невозможна, записывает в w ошибку.
func writeJson(w http.ResponseWriter, answer any) {
	resp, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(resp)
}

// writeJsonErr конвертирует переданную в err ошибку в json-формат и записывает в ответ w.
func writeJsonErr(w http.ResponseWriter, err error) {
	var jsErr JsonErr
	jsErr.Err = err.Error()
	writeJson(w, jsErr)
}
