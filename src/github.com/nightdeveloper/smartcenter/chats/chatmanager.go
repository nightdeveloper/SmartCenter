package chats

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nightdeveloper/smartcenter/settings"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type ChatManager struct {
	c *settings.Config
	bot *tgbotapi.BotAPI;
	chatChannel chan string
}

func (cm *ChatManager) Init(c *settings.Config) {
	cm.c = c

	cm.chatChannel = make(chan string)
}

func (cm *ChatManager) GetChatChannel() chan string {
	return cm.chatChannel
}

func (cm *ChatManager) startWriteLoop() {

	for msg := range cm.chatChannel {

		log.Println("sending " + msg)

		response, err := http.PostForm("http://api.pushover.net/1/messages.json", url.Values{
			"token":   {cm.c.PushToken},
			"user":    {cm.c.PushUserId},
			"device":  {cm.c.PushDeviceName},
			"title":   {"SmartCenter"},
			"message": {msg},
		})

		if err != nil {
			log.Println("error - " + err.Error())

		} else {

			defer response.Body.Close()

			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				log.Println("response read error - " + err.Error())
			}

			log.Println(fmt.Sprintf("%s\n", string(body)))
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

	// httpTransport := &http.Transport{}
	// httpClient := &http.Client{Transport: httpTransport}
	// httpTransport.Dial = dialer.Dial
	// log.Println("http transport created succesfully")


	go cm.startWriteLoop()

	log.Println("chat manager started")
}