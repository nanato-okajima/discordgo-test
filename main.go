package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type yesno struct {
	Answer string
	Forced bool
	Image  string
}

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
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
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

	// 全角？が文末にあった場合yesnoを答える
	if strings.HasSuffix(m.Content, "？") {
		url := "https://yesno.wtf/api"
		res, _ := http.Get(url)
		defer res.Body.Close()

		ba, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("レスポンスが読み込めませんでした", err)
		}
		yesno := yesno{}
		err = json.Unmarshal(ba, &yesno)
		if err != nil {
			fmt.Println("構造化に失敗しました", err)
		}
		msgimg := discordgo.MessageEmbedImage{
			URL: yesno.Image,
		}
		msgemb := discordgo.MessageEmbed{
			Title:       "結果発表",
			Image:       &msgimg,
			Description: yesno.Answer,
			Color:       100,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &msgemb)
	}

	//画像が送られて来た場合
	if len(m.Attachments) > 0 {
		for i, v := range m.Attachments {
			image := v.URL
			url := "http://whatcat.ap.mextractr.net/api_query"
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("リクエストの作成に失敗しました%#v", err)
			}
			username := os.Getenv("USER_NAME")
			password := os.Getenv("PASSWORD")
			req.SetBasicAuth(username, password)
		}
	}

	// !Helloというチャットがきたら　「Hello」　と返します
	if m.Content == "!Hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello")
	}
}
