package main

import (
	"github.com/nightdeveloper/podcastsynchronizer/rsschecker"
	psSettings "github.com/nightdeveloper/podcastsynchronizer/settings"
	"github.com/nightdeveloper/smartcenter/alivechecker"
	"github.com/nightdeveloper/smartcenter/chats"
	scSettings "github.com/nightdeveloper/smartcenter/settings"
	"io"
	"log"
	"os"
)

func initLog() {

	log.Println("hello everyone!")

	// logging
	f, err := os.OpenFile("logs/app.log", os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
	if err != nil {
		log.Fatal("error opening log file: ", err.Error())
		return
	}

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func main() {

	// init log
	initLog()

	// smart center config
	scConfig := scSettings.Config{}
	scConfig.Load()

	// podcasts synchronizer config
	psConfig := psSettings.Config{}
	psConfig.Load()

	// chat
	cm := chats.ChatManager{}
	cm.Init(&scConfig)
	cm.Start()

	chatChannel := cm.GetChatChannel()

	// alive checker loop
	ac := alivechecker.NewChecker(&scConfig)
	ac.SetChatChannel(chatChannel)
	go ac.StartLoop()

	// checker loop
	checker := rsschecker.NewChecker(&psConfig)
	checker.SetChatChannel(chatChannel)
	go checker.StartLoop()

	// infinite loop
	select{}
}