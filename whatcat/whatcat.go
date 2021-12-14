package whatcat

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	dg "github.com/bwmarrin/discordgo"
)

var url = "http://whatcat.ap.mextractr.net/api_query"

type Cat struct {
	Breed       string
	Probability float64
}

func Judge(v *dg.MessageAttachment) ([]Cat, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	_, name := filepath.Split(v.URL)
	fw, err := w.CreateFormFile("image", name)
	if err != nil {
		return nil, err
	}

	//送られてきた画像をGETで取得
	res, err := http.Get(v.URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if _, err := io.Copy(fw, res.Body); err != nil {
		return nil, err
	}
	w.Close()

	resp, err := RequestWhatCat(&b, w)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cats, err := createResult(body)
	if err != nil {
		return nil, err
	}

	return cats, nil
}

func RequestWhatCat(b *bytes.Buffer, w *multipart.Writer) (*http.Response, error) {
	//この猫なに猫APIへリクエスト準備
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return nil, err
	}

	username := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//APIからのレスポンスはJSONではなく配列を持つ配列の形式
// [["Aegean_cat", 0.735941171646],…]
func createResult(body []byte) ([]Cat, error) {
	var rs interface{}
	if err := json.Unmarshal(body, &rs); err != nil {
		return nil, err
	}

	r := rs.([]interface{})
	cats := make([]Cat, 5)
	for i, v := range r {
		cat := v.([]interface{})
		cats[i] = Cat{Breed: cat[0].(string), Probability: cat[1].(float64)}
	}

	return cats, nil
}
