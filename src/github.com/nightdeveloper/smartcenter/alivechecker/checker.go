package alivechecker;

import (
	"log"
	"time"
	"github.com/nightdeveloper/smartcenter/settings"
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
)

type Checker struct {
	config *settings.Config
	chatChannel chan string
}

func NewChecker(config *settings.Config) (c *Checker) {
	ch := new(Checker)

	ch.config = config

	return ch
}

func (c *Checker) SetChatChannel(cc chan string) {
	c.chatChannel = cc
}


func checkConnection() string {

	client := http.Client{ Timeout: time.Duration(10 * time.Minute) }
	response, err := client.Get("http://api.ipify.org/?format=json")

	if err != nil {
		log.Println("connection check error ", err.Error())

		return ""
	}

	defer response.Body.Close()

	bytes, _ := ioutil.ReadAll(response.Body)

	if err != nil {
		return ""
	}

	var ip string

	ip = string(bytes)

	ip = strings.Replace(ip, "{\"ip\":\"", "", 1)
	ip = strings.Replace(ip, "\"}", "", 1)

	return ip
}

func (c *Checker) StartLoop() {
	log.Println("alive checker loop started")

	msg := "Alive check started"
	log.Println(msg)
	c.chatChannel <- msg

	c.config.Load()

	if c.config.LastAlive != nil {
		duration := time.Since(*c.config.LastAlive);

		msg := fmt.Sprintf("Last alive at %s, unavailable %s",   
			c.config.LastAlive.Format("02.01.2006 15:04:05"),   
			duration)

		log.Println(msg)

		c.chatChannel <- msg
	}

	isConnectionNow := true
	lastIp := "";

	var lastConnectionTime time.Time

	for{
		now := time.Now();
		c.config.LastAlive = &now

		ip := checkConnection()

		isConnectionOk := ip != ""

		if isConnectionOk && lastIp == "" {
			lastIp = ip
		}

		if lastIp != ip {

			msg := fmt.Sprintf("IP changed %s -> %s",
				lastIp,
				ip)

			log.Println(msg)

			c.chatChannel <- msg

			lastIp = ip
		}

		if lastConnectionTime.IsZero() && isConnectionNow {
			lastConnectionTime = time.Now()
		}

		if isConnectionOk != isConnectionNow {
			if isConnectionOk {

				duration := ""

				if lastConnectionTime.IsZero() {
					duration = "from startup"
				} else {
					duration = time.Since(lastConnectionTime).String();
				}


				msg := fmt.Sprintf("Internet is up at %s, unavailable %s, now ip %s",
					now.Format("02.01.2006 15:04:05"),
					duration, ip)

				log.Println(msg)

				c.chatChannel <- msg
			}

			isConnectionNow = isConnectionOk
		}

		if isConnectionNow {
			lastConnectionTime = time.Now()
		}


		c.config.Save()

		time.Sleep(time.Duration(5) * time.Minute);
	}
}
