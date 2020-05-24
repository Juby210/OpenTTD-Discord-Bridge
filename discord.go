package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
)

var (
	client    *disgord.Client
	channelID disgord.Snowflake

	statsChannel chan string
	sendStats    bool
	sendStats2   bool

	clientsChannel chan string
	sendClients    bool

	cmdOutChannel chan string
	cmdOut        int
)

func discordLogin() {
	log.Println("*** Logging to discord")
	channelID = disgord.ParseSnowflakeString(config.ChannelID)

	client = disgord.New(disgord.Config{
		BotToken: config.Token,
		Logger:   disgord.DefaultLogger(false),
	})

	filter, _ := std.NewMsgFilter(ctx, client)

	client.On(disgord.EvtReady, func() {
		u, err := client.GetCurrentUser(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("*** Logged as %s", u.Tag())
	})

	client.On(disgord.EvtMessageCreate, filter.NotByBot, handleDiscordMessage)
}

func handleDiscordMessage(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message
	args := strings.Split(msg.Content, " ")
	command := strings.Replace(args[0], config.Prefix, "", 1)

	if msg.ChannelID == channelID {
		splited := strings.Split(strings.Replace(msg.Content, `"`, "'", -1), "\n")
		for _, l := range splited {
			stdin.Write([]byte(fmt.Sprintf("say \"[Discord] %s: %s\"\n", msg.Author.Username, l)))
			time.Sleep(time.Millisecond)
		}
	}

	switch command {
	case "stats":
		statsChannel = make(chan string)
		sendStats = true
		stdin.Write([]byte("server_info\n"))
		go func() {
			stats := fmt.Sprintf("```%s\n%s\n%s\n", <-statsChannel, <-statsChannel, <-statsChannel)
			sendStats = false
			sendStats2 = true
			stdin.Write([]byte("companies\n"))
			go func() {
				for sendStats2 {
					v, ok := <-statsChannel
					if ok {
						stats += "\n" + strings.Split(v, "(T:")[0]
					}
				}
			}()
			go func() {
				time.Sleep(2 * time.Millisecond)
				sendStats2 = false
				close(statsChannel)
				msg.Reply(ctx, s, stats+"```")
			}()
		}()
	case "clients":
		clientsChannel = make(chan string)
		sendClients = true
		clients := "```"
		stdin.Write([]byte("clients\n"))
		go func() {
			for sendClients {
				v, ok := <-clientsChannel
				if ok {
					clients += "\n" + strings.Split(v, "IP:")[0]
				}
			}
		}()
		go func() {
			time.Sleep(2 * time.Millisecond)
			sendClients = false
			close(clientsChannel)
			msg.Reply(ctx, s, clients+"```")
		}()
	case "help":
		msg.Reply(ctx, s, "https://github.com/Juby210/OpenTTD-Discord-Bridge#commands")
	}

	if !contains(config.Admins, msg.Author.ID.String()) {
		return
	}

	switch command {
	case "save":
		if len(args) < 2 {
			msg.Reply(ctx, s, "```Enter file name```")
			return
		}
		stdin.Write([]byte("save " + strings.Join(args[1:], " ") + "\n"))
		msg.Reply(ctx, s, "Saved")
	case "load":
		if len(args) < 2 {
			msg.Reply(ctx, s, "```Enter file name```")
			return
		}
		cmdOutChannel = make(chan string, 1)
		cmdOut = 1
		stdin.Write([]byte("load " + strings.Join(args[1:], " ") + "\n"))
		var ok bool
		go func() {
			var out string
			out, ok = <-cmdOutChannel
			if ok {
				msg.Reply(ctx, s, "```"+out+"```")
			}
		}()
		go func() {
			time.Sleep(2 * time.Millisecond)
			close(cmdOutChannel)
			if !ok {
				msg.Reply(ctx, s, "```Save loaded```")
			}
		}()
	case "restart":
		if len(args) > 1 {
			config.Args = args[1:]
		}
		cmd.Process.Kill()
	case "reset":
		if len(args) < 2 {
			msg.Reply(ctx, s, "```Enter company id```")
			return
		}
		cmdOutChannel = make(chan string, 1)
		cmdOut = 1
		stdin.Write([]byte("reset_company " + args[1] + "\n"))
		go func() {
			msg.Reply(ctx, s, "```"+<-cmdOutChannel+"```")
			close(cmdOutChannel)
		}()
	case "eval":
		if len(args) < 2 {
			msg.Reply(ctx, s, "```Enter command```")
			return
		}
		stdin.Write([]byte(strings.Join(args[1:], " ") + "\n"))
	}
}
