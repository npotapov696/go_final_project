package main

import (
	"fmt"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {

	err := db.Init()
	if err != nil {
		fmt.Printf("Ошибка подключения к БД: %s", err.Error())
		return
	}

	if err := server.RunServer(); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
