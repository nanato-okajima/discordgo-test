package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"discordgo2/whatcat"
	"discordgo2/yesno"
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
	d, err := dg.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register ready as a callback for the ready events.
	d.AddHandler(ready)

	// Register the messageCreate func as a callback for MessageCreate events.
	d.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	d.Identify.Intents = dg.MakeIntent(dg.IntentsGuilds | dg.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = d.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer d.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *dg.Session, event *dg.Ready) {
	// Set the playing status.
	log.Println("BotName: ", event.User.ID)
	log.Println("BotID: ", event.User.Username)
	s.UserUpdateStatus(dg.StatusOnline)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *dg.Session, m *dg.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	// ボットからのメッセージの場合は返さないように判定します。
	if m.Author.ID == s.State.User.ID {
		return
	}

	// !Helloというチャットがきたら　「Hello」　と返します
	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
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

	// 全角？が文末にあった場合yesnoを答える
	if strings.HasSuffix(m.Content, "？") {
		msgemb, err := yesno.Judge()
		if err != nil {
			log.Println("yesno判定エラー", err)
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &msgemb)
	}

	//画像が送られて来た場合何の猫か答える
	if len(m.Attachments) > 0 {
		answer := "この猫が・・・\n"
		for _, v := range m.Attachments {
			cats, err := whatcat.Judge(v)
			if err != nil {
				log.Println("何猫エラー", err)
			}
			for _, cat := range cats {
				answer = answer + fmt.Sprintf("%sである確率%.0f%%\n", cat.Breed, math.Floor(cat.Probability*100))
			}

			s.ChannelMessageSend(m.ChannelID, answer)
		}
	}
}
