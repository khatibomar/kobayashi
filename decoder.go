package kobayashi

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	GdriveGetDirectLinkPath = "https://www.googleapis.com/drive/v3/files/%s?fields=webContentLink"
)

var (
	ErrNotStatusOK = fmt.Errorf("Got a non %d http response", http.StatusOK)
)

type decoder interface {
	Decode(string) (string, error)
}

type Decoder struct {
}

func NewDecoder() decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(urll string) (string, error) {
	if strings.Contains(urll, "mediafire") {
		return d.mediafire(urll)
	} else if strings.Contains(urll, "drive.google") {
		return d.gdrive(urll)
	} else if strings.Contains(urll, "mixdrop") {
		return d.mixdrop(urll)
	} else if strings.Contains(urll, "fembed") {
		return d.fembed(urll)
	} else if strings.Contains(urll, "ok.ru") {
		return d.okru(urll)
	}

	return "", fmt.Errorf("host is not supported, yet")
}

func (d *Decoder) mediafire(url string) (string, error) {
	mediafireReg := `https?:\/\/download\d+.mediafire.com\/\w+\/\w+\/.*\.mp4`
	re := regexp.MustCompile(mediafireReg)
	content, status, err := httpRequest(url, http.MethodGet)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	s := re.FindString(content)
	if s == "" {
		return "", fmt.Errorf("failed to get direct link")
	}
	return s, nil
}

func (d *Decoder) gdrive(url string) (string, error) {
	var hash string
	fmt.Sscanf(url, "https://drive.google.com/file/d/%s/view?usp=sharing", &hash)
	hash = strings.TrimSuffix(hash, "/view?usp=sharing")
	u := fmt.Sprintf(GdriveGetDirectLinkPath, hash)
	content, status, err := httpRequest(u, http.MethodGet)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	return content, nil
}

func (d *Decoder) mixdrop(url string) (string, error) {
	content, status, err := httpRequest(url, http.MethodGet)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	u := NewUnpacker()
	res, err := u.Unpack(content)
	if err != nil {
		return "", err
	}
	mixdropRegx := `wurl=\"([^\"]+)`
	re := regexp.MustCompile(mixdropRegx)
	res = re.FindString(res)
	res = strings.TrimPrefix(res, `wurl="`)
	return "https:" + res, nil
}

func (d *Decoder) fembed(url string) (string, error) {
	type fembed struct {
		Data []struct {
			File string `json:"file"`
		} `json:"data"`
	}
	content, status, err := httpRequest(url, http.MethodPost)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	var fd fembed
	err = json.Unmarshal([]byte(content), &fd)
	if err != nil {
		return "", err
	}
	if len(fd.Data) == 0 {
		return "", fmt.Errorf("No direct Link Available")
	}
	return fd.Data[len(fd.Data)-1].File, nil
}

func (d *Decoder) okru(urll string) (string, error) {
	u, err := url.Parse(urll)
	if err != nil {
		return "", err
	}
	u.Host = strings.TrimPrefix(u.Hostname(), "m.")
	u.Path = strings.Replace(u.Path, "video", "videoembed", 1)
	urll = u.String()
	type okru struct {
		Flashvars struct {
			Metadata string `json:"metadata"`
		} `json:"flashvars"`
	}

	type okruMeta struct {
		Videos []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"videos"`
	}

	content, status, err := httpRequest(urll, http.MethodPost)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	regEx := "data-options=\"(.*?)\""
	re := regexp.MustCompile(regEx)
	res := re.FindString(content)
	res = strings.TrimSuffix(strings.TrimPrefix(res, `data-options="`), `"`)
	res = html.UnescapeString(res)

	var okr okru
	var okrm okruMeta

	err = json.Unmarshal([]byte(res), &okr)
	if err != nil {
		return "", err
	}
	metadata := okr.Flashvars.Metadata
	metadata = strings.TrimSuffix(strings.TrimSuffix(metadata, "{{"), "}}")
	err = json.Unmarshal([]byte(metadata), &okrm)
	if err != nil {
		return "", err
	}
	hightestURL := okrm.Videos[len(okrm.Videos)-1].URL
	return hightestURL, nil
}

func httpRequest(url, method string) (string, int, error) {
	var body io.Reader
	req, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return "", -1, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", -1, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", -1, err
	}
	defer resp.Body.Close()
	return string(b), resp.StatusCode, nil
}
