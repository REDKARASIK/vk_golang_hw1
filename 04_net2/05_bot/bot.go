package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ! ngrok deprecated !

/*
для heroku
находясь в корне репы
git subtree push --prefix 04_net2/05_bot heroku master
*/

// https://api.telegram.org/bot1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4/getUpdates

// ngrok http 8080

// ssh -R 80:localhost:8080 serveo.net

var (

	// хранить секреты в коде плохо, поэтому используем конфиги
	// "1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4"
	BotToken = flag.String("tg.token", "", "token for telegram")

	// это не секрет, но для простоты тоже выносим под конфиг
	// "https://5872-95-165-1-28.eu.ngrok.io"
	WebhookURL = flag.String("tg.webhook", "", "webhook addr for telegram")

	// запуск выглядит так:
	// ! ngrok deprecated !
	// go run bot.go -tg.token="1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4" -tg.webhook="https://5872-95-165-1-28.eu.ngrok.io"

	// serveo.net
	// go run bot.go -tg.token="1953480583:AAEU7eBaZnCUt525oUkCMRCQxK1TJmaoVd4" -tg.webhook="https://addb16f030b88ecaae1f7c65378b75b5.serveo.net"
)

var rss = map[string]string{
	"Habr": "https://habrahabr.ru/rss/best/",
}

type RSS struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	URL   string `xml:"guid"`
	Title string `xml:"title"`
}

func getNews(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	rss := new(RSS)
	err = xml.Unmarshal(body, rss)
	if err != nil {
		return nil, err
	}

	return rss, nil
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	bot, err := tgbotapi.NewBotAPI(*BotToken)
	if err != nil {
		log.Fatalf("NewBotAPI failed: %s", err)
	}

	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(*WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	updates := bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	// получаем все обновления из канала updates
	for update := range updates {
		log.Printf("upd: %#v\n", update)
		url, ok := rss[update.Message.Text]
		if !ok {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				`there is only Habr feed availible`,
			)

			msg.ReplyMarkup = &tgbotapi.ReplyKeyboardMarkup{
				Keyboard: [][]tgbotapi.KeyboardButton{
					{
						{
							Text: "Habr",
						},
					},
				},
			}
			bot.Send(msg)
			continue
		}

		rss, err := getNews(url)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"sorry, error happend",
			))
		}
		item := rss.Items[rand.Intn(len(rss.Items))]
		bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			item.URL+"\n"+item.Title,
		))
	}
}
