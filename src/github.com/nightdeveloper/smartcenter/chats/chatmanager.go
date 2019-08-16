package chats

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nightdeveloper/smartcenter/settings"
	"log"
	"net/http"
	"strings"
	"time"
)

type ChatManager struct {
	c *settings.Config
	bot *tgbotapi.BotAPI;
	chatChannel chan string
}

func (cm *ChatManager) Init(c *settings.Config) {
	cm.c = c;

	cm.chatChannel = make(chan string);
}

func (cm *ChatManager) GetChatChannel() chan string {
	return cm.chatChannel;
}

func (cm *ChatManager) startReadLoop() {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := cm.bot.GetUpdatesChan(u)

	if err != nil {
		log.Println("Can not get updates channel", err)
		return;
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if int64(update.Message.From.ID) != cm.c.TelegramOpId {
			log.Println(fmt.Sprintf("[%d %s] sends me unauth message: %s",
				update.Message.From.ID, update.Message.From.UserName, update.Message.Text))
			continue
		}

		var processed = strings.ToLower(strings.Trim(update.Message.Text, " "));

		log.Println(fmt.Sprintf("operator sends me %s, processed [%s]", update.Message.Text, processed))

		if processed == "test" {
			cm.chatChannel <- "I'm alive!"
		}
	}
}

func (cm *ChatManager) startWriteLoop() {

	for msg := range cm.chatChannel {

		newMsg := tgbotapi.NewMessage(cm.c.TelegramOpId, msg)
		newMsg.ParseMode = "markdown"

		for {
			_, err := cm.bot.Send(newMsg)

			if err != nil {
				log.Println("error while sending message: " + err.Error())
				time.Sleep(time.Duration(1) * time.Minute);
			} else {
				break;
			}
		}
	}
}

func (cm *ChatManager) Start() {

	/*
	auth := proxy.Auth{
		User: cm.c.ProxyUser,
		Password: cm.c.ProxyPass,
	}

	dialer, err := proxy.SOCKS5("tcp", cm.c.ProxyUrl, &auth, proxy.Direct)
	if err != nil {
		log.Println("can't connect to the proxy:", err)
		os.Exit(1)
	}
	*/

	log.Println("creating chat manager")

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	//httpTransport.Dial = dialer.Dial

	log.Println("http transport created succesfully")

	bot, err := tgbotapi.NewBotAPIWithClient(cm.c.TelegramKey, httpClient)

	if err != nil {
		log.Println("telegram bot creating failed", err)
		return
	}

	bot.Debug = true

	log.Println("telegram bot created")

	cm.bot = bot

	go cm.startReadLoop()
	go cm.startWriteLoop()

	log.Println("chat manager started")

}