package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	MAXLIST := 1000    // max size of list
	MAXMINUTES := 60   // max minutes to keep data
	list := []string{} // list
	lastAdd := 0       // count down from last add

	s := gocron.NewScheduler(time.UTC)
	// check every minute
	s.Cron("*/1 * * * *").Do(func() {
		//fmt.Printf("*/1 * * * *\n")
		// TODO: reset by keeper, name-pair
		if lastAdd == 0 {
			fmt.Printf("timed reset\n")
			list = []string{}
		}
		if lastAdd > 0 {
			lastAdd--
		}
	})
	s.StartAsync()

	r := gin.Default()

	// DEBUG
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// API
	// route: (url)/add?data=(data)
	// add data to list
	r.GET("/add", func(c *gin.Context) {
		data := c.DefaultQuery("data", "")
		// limit size of list
		if len(list) >= MAXLIST {
			list = list[1:]
		}
		// assume data=<keeper>,backcolor1,backcolor2,color1,color2,name1,name2,sets1,sets2,score1,score2,posession
		// ie. shannon,#000000,#ffffff,#ffffff,#000000,Them,Us,0,0,10,8,0
		// may have spaces declared as %20, looks like gin converts them to ' '

		// TODO: prefix time stamp
		currentTime := time.Now()
		entry := currentTime.Format("2006-01-02_15:04:05") + ", " + data
		list = append(list, entry)
		c.String(200, "")

		lastAdd = MAXMINUTES
	})

	// DEBUG
	// route: /raw
	// dump all list in raw form
	r.GET("/raw", func(c *gin.Context) {
		msg := strings.Join(list, "\n")
		c.String(200, msg)
	})

	// API
	// route: /reset
	// reset list
	r.GET("/reset", func(c *gin.Context) {
		list = []string{}
		c.String(200, "")
	})

	// DIRECT USER
	// route: /
	// html with list of keepers links
	r.GET("/", func(c *gin.Context) {
		msg := strings.Join(list, "\n")
		c.String(200, msg)
	})

	// DEBUG
	// route: /(keeper)/raw
	// all raw results for keeper

	// DIRECT USER
	// route: /(keeper)
	// html with results for keeper

	// API
	// route: /(keeper)/json
	// json with all results for keeper

	// API
	// route: /(keeper)/count
	// json with count of keeper results

	// API
	// route: /(keeper)/(index)
	// json with list[index] results of keeper

	r.Run(":3000")
}
