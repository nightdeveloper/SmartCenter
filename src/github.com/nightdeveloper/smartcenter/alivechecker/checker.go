package alivechecker;

import (
	"log"
	"time"
	"github.com/nightdeveloper/smartcenter/settings"
	"fmt"
)

type Checker struct {
	config *settings.Config
	chatChannel chan string
}

func NewChecker(config *settings.Config) (c *Checker) {
	ch := new(Checker);

	ch.config = config;

	return ch;
}

func (c *Checker) SetChatChannel(cc chan string) {
	c.chatChannel = cc;
}

func (c *Checker) StartLoop() {
	log.Println("alive checker loop started")

	msg := "Alive check started";
	log.Println(msg);
	c.chatChannel <- msg;

	c.config.Load()

	if (c.config.LastAlive != nil) {
		duration := time.Since(*c.config.LastAlive);

		msg := fmt.Sprintf("Last alive at %s, unavailable %02.0f:%02.0f:%02.0f",
			c.config.LastAlive.Format("02.01.2006 15:04:05"),
			duration.Hours(),
			duration.Minutes(),
			duration.Seconds());

		log.Println(msg);

		c.chatChannel <- msg
	}

	for{
		now := time.Now();
		c.config.LastAlive = &now
		c.config.Save()

		time.Sleep(time.Duration(5) * time.Minute);
	}
}
