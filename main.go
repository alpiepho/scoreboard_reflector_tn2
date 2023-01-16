package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

var VERSION string = "2.2h"

var MAXLIST int = 1000 // max size of list
var MAXLOGS int = 1000 // max size of logs

type Keeper struct {
	name string
}

var list = []string{}
var keepers = []Keeper{}

var logs = []string{}
var logOn = false

func logAdd(c *gin.Context, msg string) {
	if logOn {
		currentTime := time.Now()
		ts := currentTime.Format("2006-01-02_15:04:05")

		msg = ts + ": " + msg
		if c != nil {
			msg = ts + " (" + c.RemoteIP() + "): " + msg
		}
		if len(logs) >= MAXLOGS {
			logs = logs[:len(logs)-1] // remove last
		}
		logs = append([]string{msg}, logs...) // add first
	}
}

func getKeepersIndex(keeperName string) int {
	index := -1
	for i, keeper := range keepers {
		if keeperName == keeper.name {
			index = i
		}
	}
	return index
}

func getKeepersCount(keeper string) int {
	result := 0
	for i := 0; i < len(list); i++ {
		// ie. timestamp,keeper,...
		parts := strings.Split(list[i], ",")
		if keeper == parts[1] {
			result++
		}
	}
	return result
}

func getKeepersList(keeper string) []string {
	result := []string{}
	for i := 0; i < len(list); i++ {
		// ie. timestamp,keeper,...
		parts := strings.Split(list[i], ",")
		if keeper == "all" || keeper == "All" || keeper == parts[1] {
			result = append(result, list[i])
		}
	}
	return result
}

func getKeepersListMany(keeperNames []string) []string {
	result := []string{}
	for i := 0; i < len(list); i++ {
		// ie. timestamp,keeper,...
		parts := strings.Split(list[i], ",")

		for j := 0; j < len(keeperNames); j++ {
			keeper := keeperNames[j]
			if keeper == parts[1] || keeper == "*" {
				result = append(result, list[i])
			}
		}
	}
	return result
}

func removeKeepersList(keeper string) []string {
	result := []string{}
	for i := 0; i < len(list); i++ {
		// ie. timestamp,keeper,...
		parts := strings.Split(list[i], ",")
		if keeper != parts[1] {
			result = append(result, list[i])
		}
	}
	return result
}

// func removeKeeper(keeper string) []Keeper {
// 	result := []Keeper{}
// 	for i := 0; i < len(keepers); i++ {
// 		if keeper != keepers[i].name {
// 			result = append(result, keepers[i])
// 		}
// 	}
// 	return result
// }

const HTML_START string = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width", initial-scale=1.0"/>
    <meta name="Description" content="ScoresTN2 Reflector">
    <meta name="theme-color" content="#d36060"/>
    <title>
    ScoresTN2 Reflector
    </title>
    <link rel="stylesheet" href="/style.css" />
    <link rel="manifest" href="/manifest.json" />
    <link rel="icon"
      type="image/png" 
      href="/favicon.ico" />
  </head>
  <body class="body">
    <main>
`

const HTML_INTRO_KEEPERS string = `
    <article class="page">
      <h1  id="top">ScoresTN2 Reflector - Current Score Keepers</h1>

      <div class="introduction">
      <p>
      Version VERSION
      </p>
      <p>
      This is a list of the current score keepers.
      </p>
      </div>
`

const HTML_INTRO_ADMIN string = `
    <article class="page">
      <h1  id="top">ScoresTN2 Reflector - Current Score Keepers</h1>

      <div class="introduction">
      <p>
      Version VERSION
      </p>
      <p>
      <b>WARNING</b>: this the ADMIN page.
      </p>
      </div>
`

const HTML_INTRO_SCORES string = `
    <article class="page">
      <h1  id="top">ScoresTN2 Reflector - Scores</h1>
`

const HTML_LAST string = `
      <p>
      (This list is generated from a tool that can be found
      <a
        href="https://github.com/alpiepho/scoreboard_reflector_tn2"
        target="_blank"
        rel="noreferrer"
      >here</a>.)
      </p>
    <div id="bottom"></div>
    </article>
	</main>
  </body>
</html>
`

func buildKeepersHtml() string {
	//QUESTION: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += strings.Replace(HTML_INTRO_KEEPERS, "VERSION", VERSION, 1)
	result += "      <ul>\n"
	for _, keeper := range keepers {
		result += "        <li><a href=\"./" + keeper.name + "/html\">" + keeper.name + "</a></li>\n"
	}
	result += "        <li><a href=\"./all/html\">All</a></li>\n"
	result += "      </ul>\n"
	result += HTML_LAST
	return result
}

func buildKeeperScoresHtml(keeper string, keeperList []string) string {
	all := (keeper == "all" || keeper == "All")

	//QUESTION: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += HTML_INTRO_SCORES

	result += "      <div class=\"introduction\">\n"
	result += "      <p>\n"
	result += "      Scores for " + keeper + "\n"
	result += "      </p>\n"
	result += "      </div>\n"

	result += "      <ul>\n"

	for i := 0; i < len(keeperList); i++ {
		//fmt.Printf("%s\n", keeperList[i])
		result += "        <li>\n"
		parts := strings.Split(keeperList[i], ",")
		// time keeper colorA1 colorA2 colorB1 colorB2 nameA nameB setsA setsB scoreA scoreB possesion font zoom     sets5    setsShow
		// 0    1      2       3       4       5       6     7     8     9     10     11     12        13   14       15       16
		//                                                                                   1|2       str  zoomOn|  sets5|   setsShowOn|
		//    	                                                                                              zoomOff  sets3    setsShowOff
		if all {
			index := getKeepersIndex(parts[1])
			prefix := ""
			for i := 0; i < index; i++ {
				prefix += "&nbsp;&nbsp;"
			}
			result += prefix
		}
		if len(parts) >= 13 {
			result += parts[0] + ", "
			temp := ""
			nameA := parts[6]
			nameB := parts[7]
			colorA1 := parts[2][2:]
			colorA2 := parts[3][2:]
			colorB1 := parts[4][2:]
			colorB2 := parts[5][2:]
			if all {
				result += parts[1] + ", "
			}
			if parts[12] == "1" {
				nameA = "*" + nameA
			}
			if all {
				// color nameA
				temp := fmt.Sprintf("<span style=\"color:#%v;background-color:#%v;\">%v</span>", colorA1, colorA2, nameA)
				result += temp + ", "
			} else {
				result += nameA + ", "
			}

			if parts[12] == "2" {
				nameB = "*" + nameB
			}
			if all {
				// color nameB
				temp = fmt.Sprintf("<span style=\"color:#%v;background-color:#%v;\">%v</span>", colorB1, colorB2, nameB)
				result += temp + ", "
			} else {
				result += nameB + ", "
			}

			result += parts[10] + ", "
			result += parts[11] //+ ", "
		} else {
			// format any links
			line := keeperList[i]
			if strings.Contains(line, "https://") {
				start := strings.Index(line, "https://")
				end := strings.Index(line[start:], " ")
				if end == -1 {
					end = len(line)
				}
				url := line[start:end]
				anchor := "<a href=\"" + url + "\" target=\"_blank\" rel=\"noreferrer\">" + url + "</a>"
				line = strings.Replace(line, url, anchor, 1)
			}
			result += line
		}
		result += "        </li>\n"
	}
	result += "      </ul>\n"
	result += HTML_LAST
	return result
}

func buildAdminHtml() string {
	//QUESTION: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += strings.Replace(HTML_INTRO_ADMIN, "VERSION", VERSION, 1)
	result += "      <ul>\n"
	for _, keeper := range keepers {
		result += "        <li>\n"
		result += "          "
		result += "[<a href=\"/___/___/" + keeper.name + "/raw\">raw</a>"
		result += "&nbsp;&nbsp;"
		result += "<a href=\"/___/___/" + keeper.name + "/reset\">reset</a>"
		result += "&nbsp;&nbsp;"
		result += "<a href=\"/___/___/" + keeper.name + "/comment\">comment</a>" // TODO add comment/adjust page
		result += "]&nbsp;&nbsp;&nbsp;&nbsp;"
		result += "<a href=\"/" + keeper.name + "/html\">" + keeper.name + "</a>"
		result += "\n"
		result += "        </li>\n"
	}
	result += "      </ul>\n"
	result += "      <br>\n"
	result += "      <ul>\n"
	result += "        <li><a href=\"/___/___/raw\">All Raw</a></li>\n"
	result += "        <li><a href=\"/___/___/reset\">All Reset</a></li>\n"
	result += "        <li><a href=\"/___/___/logs\">Logs</a></li>\n"
	state := "OFF"
	if logOn {
		state = "ON"
	}
	result += "        <li><a href=\"/___/___/logs_toggle\">Logs Toggle: " + state + "</a></li>\n"
	result += "        <li><a href=\"/___/___/logs_clear\">Logs Clear</a></li>\n"
	result += "      </ul>\n"
	result += HTML_LAST
	return result
}

func buildAdminBackHtml() string {
	//QUESTION: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += strings.Replace(HTML_INTRO_ADMIN, "VERSION", VERSION, 1)
	result += "      <ul>\n"
	result += "        <li><a href=\"/___/___\">Back</a></li>\n"
	result += "      </ul>\n"
	result += HTML_LAST
	return result
}

func buildAdminCommentHtml(keeper string) string {
	//QUESTION: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += strings.Replace(HTML_INTRO_ADMIN, "VERSION", VERSION, 1)

	result += "<form method=\"POST\" action=\"/___/___/content\">\n"
	result += "    <label>Comment: </label><input name=\"comment\" type=\"text\" value=\"" + keeper + ": (from admin) \" />\n"
	result += "&nbsp;&nbsp;"
	result += "    <input type=\"submit\" value=\"Save\" />\n"
	result += "</form>\n"
	result += "<a href=\"/___/___\">Back</a></li>\n"
	result += HTML_LAST
	return result
}

func main() {
	s := gocron.NewScheduler(time.UTC)
	// every midgnight
	s.Cron("0 0 * * *").Do(func() {
		// fmt.Printf("0 0 * * *\n")
		logAdd(nil, "cron: reset all")
		list = []string{}
		keepers = []Keeper{}
	})
	// every 5 minutes
	s.Cron("*/5 * * * *").Do(func() {
		// fmt.Printf("*/5 * * * *\n")
		total := len(list)
		msg := fmt.Sprintf("cron: stats: total=%v.", total)
		logAdd(nil, msg)
	})

	s.StartAsync()

	r := gin.Default()
	r.SetTrustedProxies(nil) // clears error, but not sure if it will block when hosted

	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/style.css", "./resources/style.css")

	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	// API - for Score and Tap Applications
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////

	// API
	// route: /version
	r.GET("/version", func(c *gin.Context) {
		logAdd(c, "/version")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, VERSION)
	})

	// API
	// route: (url)/add?data=(data)
	// add data to list
	r.GET("/add", func(c *gin.Context) {
		//logAdd(c, "/add")
		data := c.DefaultQuery("data", "")
		// limit size of list
		if len(list) >= MAXLIST {
			// list = list[1:] // remove first
			list = list[:len(list)-1] // remove last
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
		keeperIndex := getKeepersIndex(parts[0])
		if keeperIndex == -1 {
			keepers = append(keepers, Keeper{
				name: parts[0],
			})
		}

		// prefix time stamp
		currentTime := time.Now()
		entry := currentTime.Format("2006-01-02_15:04:05") + "," + data
		//list = append(list, entry)              // add last
		list = append([]string{entry}, list...) // add first
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, "")

		// Build log entry
		//2023-01-15_11:19:15,dude,ff000000,fff44336,ff000000,ff448aff,Away,Home,0,0,23,14,2,FontTypes.system,zoomOff,sets3,setsShowOff
		//                  0    1        2        3        4        5    6    7 8 9 10 11 12               13     14    15          16
		parts = strings.Split(entry, ",")
		//partial := string(entry)
		partial := "/add: "
		//partial += parts[0] + ","
		partial += parts[1] + ","
		partial += "..,"
		partial += parts[6] + ","
		partial += parts[7] + ","
		partial += parts[8] + ","
		partial += parts[9] + ","
		partial += parts[10] + ","
		partial += parts[11] + ","
		partial += parts[12] + ","
		partial += "..,"
		logAdd(c, partial)
	})

	// API
	// route: /json
	// dump all list in raw form
	r.GET("/json", func(c *gin.Context) {
		logAdd(c, "/json")
		offset := 0
		count := len(list)
		matchNames := []string{"*"}

		paramPairs := c.Request.URL.Query()
		for key, values := range paramPairs {
			if key == "offset" {
				offset, _ = strconv.Atoi(values[0])
			}
			if key == "count" {
				given, _ := strconv.Atoi(values[0])
				count = given + offset
				if given == 0 || count > len(list) {
					count = len(list)
				}
			}
			if key == "names" {
				matchNames = []string{}
				parts := strings.Split(values[0], ",")
				for i := 0; i < len(parts); i++ {
					matchNames = append(matchNames, parts[i])
				}
			}
		}

		matchList := getKeepersListMany(matchNames)
		if count > len(matchList) {
			count = len(matchList)
		}
		partialList := matchList[offset:count]

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.JSON(200, gin.H{"all": partialList})
	})

	// API
	// route: /keepers/json
	// dump keeper names in json form
	r.GET("/keepers/json", func(c *gin.Context) {
		logAdd(c, "/keepers/json")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		//DEBUG
		fmt.Println(keepers)
		names := []string{}
		for _, keeper := range keepers {
			names = append(names, keeper.name)
		}
		c.JSON(200, gin.H{"keepers": names})
	})

	// API
	// route: /(keeper)/json
	// json with all results for keeper
	r.GET("/:keeperid/json", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/"+keeperid+"/json")
		keeperList := getKeepersList(keeperid)

		offset := 0
		count := len(keeperList)

		paramPairs := c.Request.URL.Query()
		for key, values := range paramPairs {
			if key == "offset" {
				offset, _ = strconv.Atoi(values[0])
			}
			if key == "count" {
				given, _ := strconv.Atoi(values[0])
				count = given + offset
				if given == 0 || count > len(keeperList) {
					count = len(keeperList)
				}
			}
		}
		partialList := keeperList[offset:count]
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.JSON(200, gin.H{keeperid: partialList})
	})

	// API - COUNT
	// route: /(keeper)/count
	// raw with count of keeper results
	r.GET("/:keeperid/count", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/"+keeperid+"/json")
		count := getKeepersCount(keeperid)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, string(rune(count)))
	})

	// API - KEEPER JSON
	// route: /(keeper)/(index)/json
	// json with list[index] results of keeper
	r.GET("/:keeperid/:indexid/json", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		indexid := c.Param("indexid")
		logAdd(c, "/"+keeperid+"/"+indexid+"/json")
		index, _ := strconv.Atoi(indexid)
		keeperList := getKeepersList(keeperid)
		if index >= 0 && index < len(keeperList) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET")
			c.JSON(200, gin.H{"entry": keeperList[index]})

		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET")
			c.String(200, "")
		}
	})

	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	// WEB - routes for direct reflector HTML
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////

	// WEB
	// route: /html
	// html with list of keepers links
	r.GET("/html", func(c *gin.Context) {
		logAdd(c, "/html")
		msg := buildKeepersHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - HOME
	// route: /
	// html with list of keepers links
	r.GET("/", func(c *gin.Context) {
		logAdd(c, "/")
		msg := buildKeepersHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - KEEPER
	// route: /(keeper)/html
	// html with results for keeper
	r.GET("/:keeperid/html", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/"+keeperid+"/html")
		keeperList := getKeepersList(keeperid)
		msg := buildKeeperScoresHtml(keeperid, keeperList)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - KEEPER
	// route: /(keeper)
	// html with results for keeper
	r.GET("/:keeperid", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/"+keeperid)
		keeperList := getKeepersList(keeperid)
		msg := buildKeeperScoresHtml(keeperid, keeperList)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	// WEB ADMIN - routes for direct reflector HTML (admin only)
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////

	// WEB - ADMIN
	// route: /___/___   // HACK: should be protected, obscure url for now
	// reset list
	r.GET("/___/___", func(c *gin.Context) {
		logAdd(c, "/___/___")
		msg := buildAdminHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN  RAW
	// route: /___/___/raw
	// dump all list in raw form
	r.GET("/___/___/raw", func(c *gin.Context) {
		logAdd(c, "/___/___/raw")
		msg := strings.Join(list, "\n")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, msg)
	})

	// WEB -ADMIN KEEPER RAW
	// route: /___/___/(keeper)/raw
	// all raw results for keeper
	r.GET("/___/___/:keeperid/raw", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/___/___/"+keeperid+"/raw")
		keeperList := getKeepersList(keeperid)
		msg := strings.Join(keeperList, "\n")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, msg)
	})

	// WEB - ADMIN - RESET ALL
	// route: /___/___/reset
	// reset list
	r.GET("/___/___/reset", func(c *gin.Context) {
		logAdd(c, "/___/___/reset")
		list = []string{}
		keepers = []Keeper{}
		msg := buildAdminBackHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN - RESET KEEPER
	// route: /(keeper)/reset
	// reset keeper
	r.GET("/___/___/:keeperid/reset", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/___/___/"+keeperid+"/reset")
		list = removeKeepersList(keeperid)
		msg := buildAdminBackHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN - RESET KEEPER
	// route: /(keeper)/comment
	// reset keeper
	r.GET("/___/___/:keeperid/comment", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		logAdd(c, "/___/___/"+keeperid+"/comment")
		msg := buildAdminCommentHtml(keeperid)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN - LOGS
	// route: /___/___/logs
	r.GET("/___/___/logs", func(c *gin.Context) {
		//logAdd(c, "/___/___/logs")
		msg := strings.Join(logs, "\n")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, msg)
	})

	// WEB - ADMIN - LOGS TOGGLE
	// route: /___/___/logs_toggle
	r.GET("/___/___/logs_toggle", func(c *gin.Context) {
		fmt.Println(logOn)
		logOn = !logOn
		//logAdd(c, "/___/___/logs_toggle")
		msg := buildAdminBackHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN - LOGS CLEAR
	// route: /___/___/logs_clear
	r.GET("/___/___/logs_clear", func(c *gin.Context) {
		logs = []string{}
		//logAdd(c, "/___/___/logs_clear")
		msg := buildAdminBackHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// WEB - ADMIN - LOGS CLEAR
	// route: /___/___/logs_clear
	r.POST("/___/___/content", func(c *gin.Context) {
		logAdd(c, "/___/___/content")
		c.Request.ParseForm()
		comment := c.Request.FormValue("comment")
		fmt.Println(comment)

		// limit size of list
		if len(list) >= MAXLIST {
			// list = list[1:] // remove first
			list = list[:len(list)-1] // remove last
		}
		// split keeper and comment
		parts := strings.Split(comment, ":")
		keeperIndex := getKeepersIndex(parts[0])
		if keeperIndex == -1 {
			keepers = append(keepers, Keeper{
				name: parts[0],
			})
		}
		// reassemble comment
		comment = parts[0] + ", " + strings.TrimSpace(parts[1])

		// prefix time stamp
		currentTime := time.Now()
		entry := currentTime.Format("2006-01-02_15:04:05") + "," + comment
		//list = append(list, entry)              // add last
		list = append([]string{entry}, list...) // add first

		msg := buildAdminBackHtml()
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	port = ":" + port

	r.Run(port)
}
