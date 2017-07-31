package chats

import (
	"github.com/nightdeveloper/smartcenter/settings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"fmt"
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

		if (update.Message == nil) {
			continue
		}

		if (int64(update.Message.From.ID) != cm.c.TelegramOpId) {
			log.Println(fmt.Sprintf("[%d %s] sends me unauth message: %s",
				update.Message.From.ID, update.Message.From.UserName, update.Message.Text))
			continue
		}

		log.Println(fmt.Sprintf("operator sends me %s", update.Message.Text))
	}
}

func (cm *ChatManager) startWriteLoop() {

	for msg := range cm.chatChannel {

		newMsg := tgbotapi.NewMessage(cm.c.TelegramOpId, msg)
		newMsg.ParseMode = "markdown"
		_, err := cm.bot.Send(newMsg)

		if (err != nil) {
			log.Println("error while sending message: " + err.Error())
		}
	}
}

func (cm *ChatManager) Start() {

	bot, err := tgbotapi.NewBotAPI(cm.c.TelegramKey)

	if err != nil {
		log.Println("telegram bot creating failed", err)
		return
	}

	log.Println("telegram bot created")

	cm.bot = bot

	go cm.startReadLoop()
	go cm.startWriteLoop()

	log.Println("chat manager started")

}