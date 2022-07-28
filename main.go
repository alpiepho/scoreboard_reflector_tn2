package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// [ ] - Simple site that tracks the most recent scores
// [x] - Use GET with params as *input*, not real api
// [x] - Clears data at certain size
// [ ] - Clears data after certain idle time
// [ ] - Web endpoints:
// [x]     - "/" get raw text dump of all data
// [ ]     - "/scorer" get raw text data for "scorer"
// [ ]     - "/scorer/count" get count of packets for scorer, returns json
// [ ]     - "/scorer/n" get nth packet, returns json
// [x]     - "/add?s=scorer&d={comma separated data}" how to add data with GET
// [x]     - "/reset" reset application
// [ ] - data added with scorer id
// [ ] - timestamp added to data packet as received

func main() {
	MAXLIST := 1000
	list := []string{}
	clearTime := time.Minute * 3

	var timer = time.AfterFunc(clearTime, func() {
		// reset list after 30 minutes
		list = []string{}
		fmt.Printf("reset list")
	})

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
		if len(list) >= MAXLIST {
			list = []string{}
		}
		// assume data=<keeper>,backcolor1,backcolor2,color1,color2,name1,name2,sets1,sets2,score1,score2,posession
		// ie. shannon,#000000,#ffffff,#ffffff,#000000,Them,Us,0,0,10,8,0
		// may have spaces declared as %20, looks like gin converts them to ' '

		// TODO: prefix time stamp
		list = append(list, data)
		c.String(200, "")

		if timer.Stop() {
			<-timer.C
		}
		timer.Reset(clearTime)
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
