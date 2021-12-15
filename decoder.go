package kobayashi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	}
	if strings.Contains(urll, "drive.google") {
		return d.gdrive(urll)
	}
	if strings.Contains(urll, "mixdrop") {
		return d.mixdrop(urll)
	}
	if strings.Contains(urll, "fembed") {
		return d.fembed(urll)
	}
	return "", fmt.Errorf("host is not supported, yet...")
}

func (d *Decoder) mediafire(url string) (string, error) {
	mediafireReg := `https:\/\/download\d+.mediafire.com\/\w+\/\w+\/.*\.mp4`
	re := regexp.MustCompile(mediafireReg)
	content, status, err := httpRequest(url, http.MethodGet)
	if err != nil {
		return "", nil
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	return re.FindString(content), nil
}

func (d *Decoder) gdrive(url string) (string, error) {
	var hash string
	fmt.Sscanf(url, "https://drive.google.com/file/d/%s/view?usp=sharing", &hash)
	hash = strings.TrimSuffix(hash, "/view?usp=sharing")
	u := fmt.Sprintf(GdriveGetDirectLinkPath, hash)
	content, status, err := httpRequest(u, http.MethodGet)
	if err != nil {
		return "", nil
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	return content, nil
}

func (d *Decoder) mixdrop(url string) (string, error) {
	content, status, err := httpRequest(url, http.MethodGet)
	if err != nil {
		return "", nil
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
		return "", nil
	}
	if status != http.StatusOK {
		return "", ErrNotStatusOK
	}
	var fd fembed
	err = json.Unmarshal([]byte(content), &fd)
	if err != nil {
		return "", nil
	}
	if len(fd.Data) == 0 {
		return "", fmt.Errorf("No direct Link Available")
	}
	return fd.Data[len(fd.Data)-1].File, nil
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
		return "", -1, nil
	}
	defer resp.Body.Close()
	return string(b), resp.StatusCode, nil
}
