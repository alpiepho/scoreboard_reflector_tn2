package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

type keeper struct {
	name    string
	lastAdd int
}

func main() {
	MAXLIST := 1000    // max size of list
	MAXMINUTES := 60   // max minutes to keep data
	list := []string{} // list
	lastAdd := 0       // count down from last add
	keepers := []keeper{}

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
	r.SetTrustedProxies(nil) // clears error, but not sure if it will block when hosted

	// DEBUG
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// API
	// route: /reset
	// reset list
	r.GET("/reset", func(c *gin.Context) {
		list = []string{}
		c.String(200, "")
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
		// *color*: 6 char hex, rgb
		// name*: characters, use %20 for space
		// sets*, score*: integer
		// posession: 0-none, 1-them, 2-us
		// ie. shannon,000000,ffffff,ffffff,000000,Them,Us,0,0,10,8,0
		// ie. shannon,Them,Us,0,0,10,8,0
		// may have spaces declared as %20, looks like gin converts them to ' '

		parts := strings.Split(data, ",")
		if len(parts) == 12 {
			fmt.Printf("long form\n")
		}
		if len(parts) == 8 {
			fmt.Printf("long form\n")
		}
		if findKeepersIndex(parts[0], keepers) != -1 {
			fmt.Printf("keeper '%s' found\n", parts[0])
		} else {
			k := keeper{
				name:    "shannon",
				lastAdd: MAXMINUTES,
			}
			keepers = append(keepers, k)
		}

		// TODO: prefix time stamp
		currentTime := time.Now()
		entry := currentTime.Format("2006-01-02_15:04:05") + "," + data
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

	// DEBUG
	// route: /json
	// dump all list in raw form
	r.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{"all": list})
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

func findKeepersIndex(keeper string, keepers []keeper) int {
	index := -1
	for i := 0; i < len(keepers); i++ {
		if keeper == keepers[i].name {
			index = i
		}
	}
	return index
}
