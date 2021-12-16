package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"discordgo2/sendgrid"
	"discordgo2/whatcat"
	"discordgo2/yesno"
)

type translated struct {
	Text string
	code int
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
	d, err := dg.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register ready as a callback for the ready events.
	d.AddHandler(ready)

	// Register the messageCreate func as a callback for MessageCreate events.
	d.AddHandler(messageCreate)

	d.AddHandler(messageReactionAdd)

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

	// !Helloというチャットがきたら　メッセージに絵文字をつけて
	//「Hello」　と返します
	if m.Content == "!Hello" {
		err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, "👺")
		if err != nil {
			fmt.Println("リアクションに失敗しました", err)
		}
		_, err = s.ChannelMessageSend(m.ChannelID, "Hello")
		if err != nil {
			fmt.Println("Helloに失敗しました", err)
		}
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
				func() {
					url := fmt.Sprintf("https://script.google.com/macros/s/AKfycbyQvThz03giX6sSV9jZHCudENQhUYnfOimZzwhvgygbVnWyhCOZEWSYJjx5UNylbWo9Wg/exec?text=%s&source=en&target=ja", cat.Breed)

					res, err := http.Get(url)
					if err != nil {
						fmt.Println(err)
					}
					defer res.Body.Close()

					b, _ := ioutil.ReadAll(res.Body)
					tr := new(translated)
					json.Unmarshal(b, &tr)

					answer = answer + fmt.Sprintf("%sである確率%.0f%%\n", tr.Text, math.Floor(cat.Probability*100))
				}()
			}

			s.ChannelMessageSend(m.ChannelID, answer)
		}
	}

	// sendgridでメールを送る
	if m.Content == "mail" {
		sendgrid.SendMail()
		s.ChannelMessageSend(m.ChannelID, "メールを送信しました")
	}
}

func messageReactionAdd(s *dg.Session, m *dg.MessageReactionAdd) {
	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println("チャンネルメッセージの取得に失敗しました", err)
		return
	}
	usr, err := s.User(m.UserID)
	if err != nil {
		fmt.Println("ユーザーが取得できませんでした", err)
	}
	message := fmt.Sprintf("%sが%sをチェックしました。", usr.Username, msg.Content)
	s.ChannelMessageSend(m.ChannelID, message)
}
