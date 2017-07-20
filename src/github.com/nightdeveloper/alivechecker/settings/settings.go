package settings

import (
	"encoding/json"
	"log"
	"path/filepath"
	"io/ioutil"
	"time"
)

type Config struct {
	LastAlive	*time.Time	`json:"lastAlive,omitempty"`
}

func (c *Config) getFileName() string {
	absPath, _ := filepath.Abs("./");
	return absPath + "config_ac.json";
}

func (c *Config) Load() {

	file, err := ioutil.ReadFile(c.getFileName())

	if err != nil {
		log.Fatal("Config reading error from " + c.getFileName() + " ", err);
		panic("config reading error");
	}

	err = json.Unmarshal(file, c);

	if err != nil || c == nil {
		log.Fatal("Config decoding error ", err);
		panic("config decoding error");
	}

	log.Println("config read");
}

func (c *Config) Save() {

	out, err := json.MarshalIndent(c, "", "	")

	err = ioutil.WriteFile(c.getFileName(), out,0755)

	if err != nil {
		log.Fatal("Config write error ", err);
		panic("config write error");
	}

	log.Println("Config saved");
}