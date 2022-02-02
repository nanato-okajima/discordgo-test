package sendgrid

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail() {
	eapiKey := os.Getenv("API_KEY")
	etos := strings.Split(os.Getenv("TOS"), ",")
	efrom := os.Getenv("FROM")

	message := mail.NewV3Mail()
	from := mail.NewEmail("", efrom)
	message.SetFrom(from)

	//1つ目の宛先と、対応するSubstitutionタグを指定
	p := mail.NewPersonalization()
	to := mail.NewEmail("", etos[0])
	p.AddTos(to)
	p.SetSubstitution("%fullname%", "田中 太郎")
	p.SetSubstitution("%familyname%", "田中")
	p.SetSubstitution("%place%", "中野")
	message.AddPersonalizations(p)

	//件名を設定
	message.Subject = "[sendgrid-example] フクロウのお名前は%fullname%さん"
	//テキストパートを設定
	c := mail.NewContent("text/plain", "%familyname% さんは何をしていますか？\r\n 彼は%place%にいます。")
	message.AddContent(c)
	//HTMLパートを設定
	c = mail.NewContent("text/html", "<strong> %familyname% さんは何をしていますしていますか？</strong><br>彼は%place%にいます。")
	message.AddContent(c)

	// カテゴリ情報を付加
	message.AddCategories("category1")
	// カスタムヘッダを指定
	message.SetHeader("X-Sent-Using", "SendGrid-API")
	//画像ファイルを添付
	a := mail.NewAttachment()
	file, _ := os.OpenFile("./neko.jpg", os.O_RDONLY, 0666)
	defer file.Close()

	data, _ := ioutil.ReadAll(file)
	data_enc := base64.StdEncoding.EncodeToString(data)
	a.SetContent(data_enc)
	a.SetType("image/jpg")
	a.SetFilename("nekoneko.jpg")
	a.SetDisposition("attachment")
	message.AddAttachment(a)

	//メール送信を行い、レスポンスを表示
	client := sg.NewSendClient(eapiKey)
	res, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(res.StatusCode)
		fmt.Println(res.Body)
		fmt.Println(res.Headers)
	}
}
