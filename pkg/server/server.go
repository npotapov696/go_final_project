// Пакет server реализует запуск сервера.
package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go1f/pkg/api"
)

// DefaultPort содержит значение порта сервера по умолчанию.
var DefaultPort = 7540

var envPort = os.Getenv("TODO_PORT") // Получаем переменную окружения TODO_PORT.

// getPort возвращает значение порта в виде строки ":<значение порта>" из переменной среды окружения TODO_PORT.
// Если таковая отсутствует, возвращает значение пременной DefaultPort.
func getPort() int {
	port := DefaultPort
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	return port
}

// RunServer запускает сервер.
func RunServer() error {
	api.Init()

	port := getPort()
	fmt.Printf("Приложение запущено на порту: %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
