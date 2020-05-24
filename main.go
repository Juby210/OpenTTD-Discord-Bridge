package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
)

type cfg struct {
	OpenTTD, Token, ChannelID, Prefix string
	Args, Admins                      []string
}

var (
	config cfg
	ctx    = context.Background()
)

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln(err)
	}

	discordLogin()
	startOpenTTD()

	client.StayConnectedUntilInterrupted(ctx)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
