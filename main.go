package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"unicode"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	BOT_API_TOKEN = "BOT_API_TOKEN"
	BASE_URL      = "BASE_URL"
	IS_DEBUG      = "IS_DEBUG"
	PORT          = "PORT"
	TEMIR_URL     = "https://raw.githubusercontent.com/aysuegr/polytopia/master/pl_PL.json"
)

var (
	token   = os.Getenv(BOT_API_TOKEN)
	baseUrl = os.Getenv(BASE_URL)
	isDebug = os.Getenv(IS_DEBUG) != "0"
	port    = os.Getenv(PORT)
)

func main() {

	var poller tb.Poller
	var webhook *tb.Webhook

	if isDebug {
		log.Print("using polling")
		poller = &tb.LongPoller{
			Timeout: 10 * time.Second,
		}
	} else {
		log.Print("using webhooks")
		webhook = &tb.Webhook{
			Listen: ":" + port,
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: baseUrl,
			},
		}
		poller = webhook
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: poller,
	})
	if err != nil {
		errPrint(fmt.Errorf("cannot initiate bot client: %w", err))
		panic(err)
	}

	log.Print("getting webhook")
	tmpwh, err := bot.GetWebhook()
	if err != nil {
		log.Print(err)
	}

	if tmpwh != nil {
		log.Print("webhook exists, trying to delete it")
		err = bot.RemoveWebhook()
		if err != nil {
			errPrint(fmt.Errorf("cannot delete the webhook: %w", err))
		}
	}

	b := Bot{b: bot}

	bot.Handle("/esep", b.measureHandler)

	log.Print("attempting to start")
	bot.Start()
}

type Bot struct {
	b *tb.Bot
}

func errPrint(s error) {
	log.Printf("[ERROR] %v\n", s)
}

func (bot *Bot) measureHandler(m *tb.Message) {
	var msg string
	tred, cnt, err := measure()
	if err != nil {
		msg = fmt.Sprintf("аударманы хисаптағанда қате кетті: %v", err)
	} else {
		msg = fmt.Sprintf("сөздіктің %d%% аударылған, яғни %d/%d", int(float64(tred)/float64(cnt)*100),
			tred, cnt)
	}

	_, err = bot.b.Send(m.Chat, msg, tb.ModeMarkdown)
	if err != nil {
		errPrint(fmt.Errorf("cannot send message: %w", err))
	}
}

func measure() (uint64, uint64, error) {
	type translation map[string]string

	var t translation

	r, err := http.Get(TEMIR_URL)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot get data from url {%s}: %w", TEMIR_URL, err)
	}

	if err = json.NewDecoder(r.Body).Decode(&t); err != nil {
		panic(err)
	}

	var tred, cnt uint64

	for _, v := range t {
		for _, r := range v {
			if unicode.Is(unicode.Cyrillic, r) {
				tred++
				break
			}
		}
		cnt++
	}

	return tred, cnt, nil
}
