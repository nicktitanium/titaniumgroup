package main

import (
	"cryptoChange/pkg/bot"
	"cryptoChange/pkg/database"
	"log"
)

func main() {
	//token := os.Getenv("")
	token := ""
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	if err := database.InitDB("orders.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	bot.Start(token)
}
