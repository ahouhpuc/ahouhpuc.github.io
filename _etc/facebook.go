package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var facebookAccessToken = os.Getenv("AOP_FACEBOOK_ACCESS_TOKEN")

var aop Page

type JsonTime time.Time

type FacebookObject interface {
	url() string
}

type Photo struct {
	Id        string
	Source    string
	cacheUrl  string `json:"-"`
	cacheData []byte `json:"-"`
}

type Post struct {
	Id          string
	CreatedTime JsonTime `json:"created_time"`
	Message     string
	ObjectId    string `json:"object_id"`
	Type        string
	Photo       Photo
}

type Page struct {
	Posts    []Post    `json:"data"`
	cachedAt time.Time `json:"-"`
}

func (page Page) url() string {
	return "ahouhpuc/posts"
}

func (post Post) url() string {
	return post.ObjectId
}

func (photo Photo) url() string {
	return photo.Id
}

func (photo *Photo) CacheUrl() string {
	if photo.cacheUrl == "" {
		u, err := url.Parse(photo.Source)
		if err == nil {
			photo.cacheUrl = fmt.Sprintf("/photo/%s%s", photo.Id, path.Ext(u.Path))
		}
	}

	return photo.cacheUrl
}

func (photo *Photo) Data() []byte {
	if len(photo.cacheData) == 0 {
		photo.cacheData = get(photo.Source)
	}

	return photo.cacheData
}

func (t JsonTime) String() string {
	r := strings.NewReplacer(
		"January",
		"janvier",
		"February",
		"février",
		"March",
		"mars",
		"April",
		"avril",
		"May",
		"mai",
		"June",
		"juin",
		"July",
		"juillet",
		"August",
		"août",
		"September",
		"septembre",
		"October",
		"octobre",
		"November",
		"novembre",
		"December",
		"décembre",
	)

	return r.Replace(time.Time(t).Format("2 January à 15:04"))
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	u, err := time.Parse(`"2006-01-02T15:04:05-0700"`, string(data))
	*t = JsonTime(u)
	return
}

func fetch(object FacebookObject) {
	body := get(fmt.Sprintf(
		"https://graph.facebook.com/%s?access_token=%s",
		object.url(),
		facebookAccessToken,
	))

	json.Unmarshal(body, object)
}

func get(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil
	}

	return body
}

func (page *Page) cache() {
	if page.cachedAt.Add(time.Hour).Before(time.Now()) {
		fetch(page)

		ch := make(chan Photo)
		n := 0

		for i := range page.Posts {
			post := &page.Posts[i]
			if post.Type == "photo" {
				go func(p *Post) {
					post.Photo = Photo{Id: post.ObjectId}
					fetch(&(post.Photo))
					ch <- p.Photo
				}(post)
				n += 1
			}
		}

		for i := 0; i < n; i++ {
			<-ch
		}

		page.cachedAt = time.Now()
	}
}

func facebookHandler(w http.ResponseWriter, r *http.Request) {
	aop.cache()

	t, err := template.ParseFiles("_site/fr.html", "_site/facebook.html")

	if err != nil {
		panic(err)
	}

	err = t.Execute(w, aop)

	if err != nil {
		panic(err)
	}
}

func photoHandler(w http.ResponseWriter, r *http.Request) {
	for i := range aop.Posts {
		post := &aop.Posts[i]
		if post.Type == "photo" && post.Photo.cacheUrl == r.URL.Path {
			w.Write(post.Photo.Data())
			return
		}
	}

	http.NotFound(w, r)
}
