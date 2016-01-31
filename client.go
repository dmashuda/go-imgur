package imgur

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type imgurClient struct {
	clientId        string
	baseUrl         string
	httpClient      *http.Client
	userLimit       int
	userRemaining   int
	userReset       time.Time
	clientLimit     int
	clientRemaining int
}

type Comment struct {
	Id         int       `json:"id"`
	ImageId    string    `json:"image_id"`
	Comment    string    `json:"comment"`
	Author     string    `json:"author"`
	AuthorId   int       `json:"author_id"`
	OnAlbum    bool      `json:"on_album"`
	AlbumCover string    `json:"album_cover"`
	DateTime   int       `json:"datetime"`
	Ups        int       `json:"ups"`
	Downs      int       `json:"downs"`
	Points     int       `json:"points"`
	ParentId   int       `json:"parent_id"`
	Deleted    bool      `json:"deleted"`
	Children   []Comment `json:"children"`
}

type Image struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    int    `json:"datetime"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Size        int    `json:"size"`
	Views       int    `json:"views"`
	Bandwidth   int    `json:"bandwidth"`
	Link        string `json:"link"`
	Nsfw        bool   `json:"nsfw"`
	Section     string `json:"section"`
}

type ImgurGallery struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	DateTime       int       `json:"datetime"`
	Cover          string    `json:"cover"`
	Nsfw           bool      `json:"nsfw"`
	CommentCount   int       `json:"comment_count"`
	CommentPreview []Comment `json:"comment_preview"`
	Topic          string    `json:"topic"`
	ImageCount     int       `json:"images_count"`
	Images         []Image   `json:"images"`
}

const baseImgurUrl = "https://api.imgur.com/3"

func NewClient(clientId string) imgurClient {
	iClient := imgurClient{
		clientId:   clientId,
		httpClient: &http.Client{},
		baseUrl:    baseImgurUrl,
	}

	return iClient
}

func (c imgurClient) get(url string, params map[string]string, r io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("GET", c.baseUrl+url, r)
	request.Header.Add("Authorization", "Client-ID "+c.clientId)
	if err != nil {
		return nil, err
	}

	values := request.URL.Query()
	for k, s := range params {
		values.Add(k, s)
	}
	request.URL.RawQuery = values.Encode()
	resp, err := c.httpClient.Do(request)

	c.userLimit = strconv.Atoi(resp.Header.Get("X-RateLimit-UserLimit"))
	c.userRemaining = strconv.Atoi(resp.Header.Get("X-RateLimit-UserRemaining"))
	c.userReset = strconv.Atoi(resp.Header.Get("X-RateLimit-UserReset"))
	c.clientLimit, _ = time.Parse(time.UnixDate, resp.Header.Get("X-RateLimit-ClientLimit"))
	c.clientRemaining = strconv.Atoi(resp.Header.Get("X-RateLimit-ClientRemaining"))
	log.Println(c)
	return resp, err
}

func (c imgurClient) GetAlbum(url string, page, perPage int) ([]Image, error) {
	images := struct {
		Data    []Image `json:"data"`
		Success bool    `json:"success"`
		Status  int     `json:"status"`
	}{}

	params := make(map[string]string)

	params["page"] = strconv.Itoa(page)
	params["perPage"] = strconv.Itoa(perPage)

	resp, err := c.get(url, params, nil)

	if err != nil {
		return images.Data, err
	}
	respBytes, readErr := ioutil.ReadAll(resp.Body)

	if err != readErr {
		return images.Data, readErr
	}
	marshalErr := json.Unmarshal(respBytes, &images)

	return images.Data, marshalErr
}
