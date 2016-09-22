package main

import (
	"encoding/json"
	"flag"
	log "fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Name      string
	PostCount int    `json:"posts-total"`
	PostList  []Post `json:"posts"`
}

type Post struct {
	ID        json.Number `json:"id"`
	URL       string      `json:"url"`
	Type      string      `json:"type"`
	Slug      string      `json:"slug"`
	Timestamp int64       `json:"unix-timestamp"`
}

func TrimJS(c []byte) []byte {
	// The length of "var tumblr_api_read = " is 22.
	return c[22 : len(c)-2]
}

// var contents = `{"tumblelog":{"title":"Untitled","description":"","name":"reg4net","timezone":"AsiaTokyo","cname":false,"feeds":[]},"posts-start":0,"posts-total":7}`
// var contents = `{"postCount":1,"posts-total":7,"posts":[{"id":"148619653632","url":"http:\/\/reg4net.tumblr.com\/post\/148619653632"}]};`

func getBlogData(u User) {
	url := log.Sprintf("http://%s.tumblr.com/api/read/json", u.Name)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("--->>> error, can't get user:%s json.", u.Name)
		return
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("--->>> http get fail.")
	}

	contents = TrimJS(contents)

	// log.Printf("%s \n", contents)

	var blog User
	erro := json.Unmarshal([]byte(contents), &blog)
	if erro != nil {
		log.Println("--->>> json unmarshal fail.")
	}
	// log.Println(blog)
	log.Printf("count:%d\n", blog.PostCount)
	for _, v := range blog.PostList {
		log.Printf("id:%s, type:%s, url:%s \n", v.ID, v.Type, v.URL)
	}
	// ioutil.WriteFile("result.json", []byte(blog), 0644)
}

func main() {
	var u User
	flag.Parse()

	users := flag.Args()
	for _, user := range users {
		log.Println(user)
		u.PostCount = 0
		u.Name = user
		getBlogData(u)
	}

}
