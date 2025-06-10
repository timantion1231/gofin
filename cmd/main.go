package main

import (
	"financeAppAPI/internal/app"
	"log"
)

func main() {
	if err := app.InitApp(); err != nil {
		log.Fatal("Ошибка запуска приложения: ", err)
	}
}
