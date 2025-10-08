package main

// сюда писать код

import (
	"context"
)

var (
	// @BotFather в телеграме даст вам токен. Если захотите потыкать своего бота через телегу - используйте именно его
	BotToken = "XXX"

	// Урл, в который будет стучаться телега при получении сообщения от пользователя.
	// Может быть как айпишником личной виртуалки, так и просто выдан сервисом для деплоя
	WebhookURL = "https://525f2cb5.ngrok.io"
)

func startTaskBot(ctx context.Context) error {
	// Сюда пишите ваш код
	return nil
}

func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
