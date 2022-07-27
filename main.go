package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// [ ] - Simple site that tracks the most recent scores
// [ ] - Use GET with params as *input*, not real api
// [ ] - Clears data at certain size
// [ ] - Clears data after certain idle time
// [ ] - Web endpoints:
// [x]     - "/" get raw text dump of all data
// [ ]     - "/scorer" get raw text data for "scorer"
// [ ]     - "/scorer/count" get count of packets for scorer, returns json
// [ ]     - "/scorer/n" get nth packet, returns json
// [ ]     - "/add?s=scorer&d={comma separated data}" how to add data with GET
// [x]     - "/reset" reset application
// [ ] - data added with scorer id
// [ ] - timestamp added to data packet as received

func main() {
	list := []string{}
	count := 0

	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		count++
		list = append(list, fmt.Sprint(count))
		c.String(200, "Hello, World!")
	})

	r.GET("/", func(c *gin.Context) {
		msg := strings.Join(list, "\n")
		c.String(200, msg)
	})

	r.GET("/reset", func(c *gin.Context) {
		list = []string{}
		c.String(200, "")
	})

	r.Run(":3000")
}
