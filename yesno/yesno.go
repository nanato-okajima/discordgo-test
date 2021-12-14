package yesno

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	dg "github.com/bwmarrin/discordgo"
)

const url = "https://yesno.wtf/api"

type answer struct {
	Answer string
	Forced bool
	Image  string
}

func Judge() (dg.MessageEmbed, error) {
	res, err := http.Get(url)
	if err != nil {
		return dg.MessageEmbed{}, err
	}
	defer res.Body.Close()

	ba, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dg.MessageEmbed{}, err
	}

	answer := new(answer)
	err = json.Unmarshal(ba, &answer)
	if err != nil {
		return dg.MessageEmbed{}, err
	}

	msgimg := dg.MessageEmbedImage{
		URL: answer.Image,
	}
	msgemb := dg.MessageEmbed{
		Title:       "結果発表",
		Image:       &msgimg,
		Description: answer.Answer,
		Color:       100,
	}

	return msgemb, nil
}
