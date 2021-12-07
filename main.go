package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	/*local only code */
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		// .env読めなかった場合の処理
		log.Fatal(err)
	}

	Token := os.Getenv("DISCORD_TOKEN")
	log.Println("Token: ", Token)
	if Token == "" {
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	log.Println("BotName: ", event.User.ID)
	log.Println("BotID: ", event.User.Username)
	s.UserUpdateStatus(discordgo.StatusOnline)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	// ボットからのメッセージの場合は返さないように判定します。
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Server名を取得して返します。
	if m.Content == "ServerName" {
		g, err := s.Guild(m.GuildID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(g.Name)
		s.ChannelMessageSend(m.ChannelID, g.Name)
	}

	// !Helloというチャットがきたら　「Hello」　と返します
	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
	}
}
