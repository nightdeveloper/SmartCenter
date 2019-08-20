package alivechecker;

import (
	"log"
	"time"
	"github.com/nightdeveloper/smartcenter/settings"
	"fmt"
	"net/http"
	"io/ioutil"
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


func checkConnection(url string) string {

	client := http.Client{ Timeout: time.Duration(10 * time.Minute) }
	response, err := client.Get(url)

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

	return ip
}

func (c *Checker) StartLoop() {
	log.Println("Alive check started")

	c.config.Load()

	if c.config.LastAlive != nil {
		duration := time.Since(*c.config.LastAlive).Round(time.Second)

		msg := fmt.Sprintf("Restarted, last alive at %s, unavailable %s",
			c.config.LastAlive.Format("02.01.2006 15:04:05"),   
			duration)

		log.Println(msg)

		c.chatChannel <- msg
	}

	isConnectionNow := true
	lastIp := ""

	var lastConnectionTime time.Time

	for{
		now := time.Now();
		c.config.LastAlive = &now

		checkService := "1"
		ip := checkConnection(c.config.GetIPURL1)

		if ip == "" {
			checkService = ""
		}

		isConnectionOk := ip != ""

		if isConnectionOk && lastIp == "initial" {
			lastIp = ip

    			msg := fmt.Sprintf("started with ip %s (checkservice %s)", lastIp, checkService)

			log.Println(msg)

			c.chatChannel <- msg
		}

		if lastIp != ip && ip != "" {

			msg := ""

			if lastIp == "" {
				msg = fmt.Sprintf("Started with IP %s (checkservice %s)",
					ip,
					checkService)

			} else {
				msg = fmt.Sprintf("IP changed %s -> %s (checkservice %s)",
					lastIp,
					ip,
					checkService)
			}

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
					duration = time.Since(lastConnectionTime).Round(time.Second).String()
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

		time.Sleep(time.Duration(1) * time.Minute);
	}
}
