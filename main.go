package main

import (
	"encoding/json"
	"log"
	"os"
	"plexwebhooks-rich-presence/entities"
	"sync"
	"time"

	"github.com/ananagame/rich-go/client"
	"github.com/gin-gonic/gin"
)

var id string
var dc discord

type mydiscord struct {
	lastUpdate time.Time
	logined    bool
	mu         sync.Mutex
}

func main() {
	id = os.Getenv("DISCORD_APP_ID")
	if id == "" {
		log.Panic("env DISCORD_APP_ID is null")
	}

	n := newDiscord()
	dc = n

	go func() {
		for {
			if dc.isLogin() && dc.isOld() {
				dc.logout()
			}
			time.Sleep(time.Second * 10)
		}
	}()

	r := gin.Default()
	r.POST("/", plexWebHook)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func plexWebHook(c *gin.Context) {
	// [TODO]: handle thumbnail image
	// c.Request.FormFile("thumb")
	var nowPlaying entities.PlexWebHooks
	jsonStr := c.Request.FormValue("payload")
	err := json.Unmarshal([]byte(jsonStr), &nowPlaying)
	if err != nil {
		panic(err)
	}

	if !dc.isLogin() {
		dc.login()
	}

	if (nowPlaying.Event == "media.play" || nowPlaying.Event == "media.resume") && nowPlaying.Metadata.Type == "track" {
		err := client.SetActivity(client.Activity{
			State:   nowPlaying.Metadata.Title,
			Details: nowPlaying.Metadata.ParentTitle,
		})
		if err != nil {
			panic(err)
		}
	}
}

func (d *mydiscord) logout() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.logined = false
	client.Logout()
}

func (d *mydiscord) login() {
	d.mu.Lock()
	defer d.mu.Unlock()
	err := client.Login(id)
	if err != nil {
		panic(err)
	}
	d.logined = true
}

func (d *mydiscord) isLogin() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.logined
}

func (d *mydiscord) isOld() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.lastUpdate.Add(time.Minute*10).UnixNano() < time.Now().UnixNano()
}

func newDiscord() *mydiscord {
	return &mydiscord{
		lastUpdate: time.Now(),
		logined:    false,
	}
}

type discord interface {
	login()
	logout()
	isLogin() bool
	isOld() bool
}
