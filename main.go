package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adalessa/bot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	Config *config.Config
)

func init() {
	var err error
	Config, err = config.NewConfig()
	if err != nil {
		panic(err)
	}

}

func main() {
	// seed del random
	rand.Seed(time.Now().UTC().UnixNano())

	dg, err := discordgo.New("Bot " + Config.DiscordToken)

	if err != nil {
		fmt.Println("error creating Discord session", err)
		return
	}

	api := NewOpApi(Config.OpApiHost)
	listeners := NewListeners(api)

	dg.AddHandler(listeners.MessageListener())

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running, Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
