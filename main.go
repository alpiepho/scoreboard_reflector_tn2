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

type Keeper struct {
	name    string
	lastAdd int
}

func getKeepersIndex(keeper string, keepers []Keeper) int {
	index := -1
	for i := 0; i < len(keepers); i++ {
		if keeper == keepers[i].name {
			index = i
		}
	}
	return index
}

func getKeepersCount(keeper string, list []string) int {
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

func getKeepersList(keeper string, list []string) []string {
	result := []string{}
	for i := 0; i < len(list); i++ {
		// ie. timestamp,keeper,...
		parts := strings.Split(list[i], ",")
		if keeper == parts[1] {
			result = append(result, list[i])
		}
	}
	return result
}

func getKeepersListMany(keeperNames []string, list []string) []string {
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

func removeKeepersList(keeper string, list []string) []string {
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

// func removeKeeper(keeper string, keepers []Keeper) []Keeper {
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
      Version 2.2
      </p>
      <p>
      This is a list of the current score keepers.
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

func buildKeepersHtml(keepers []Keeper) string {
	//TODO: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += HTML_INTRO_KEEPERS
	result += "      <ul>\n"
	for i := 0; i < len(keepers); i++ {
		result += "        <li><a href=\"./" + keepers[i].name + "/html\">" + keepers[i].name + "</a>\n"
		result += "        </li>\n"
	}
	result += "      </ul>\n"
	result += HTML_LAST
	return result
}

func buildKeeperScoresHtml(keeper string, list []string) string {
	//TODO: consider using template file and r.HTML
	result := ""
	result += HTML_START
	result += HTML_INTRO_SCORES

	result += "      <div class=\"introduction\">\n"
	result += "      <p>\n"
	result += "      Scores for " + keeper + "\n"
	result += "      </p>\n"
	result += "      </div>\n"

	result += "      <ul>\n"

	for i := 0; i < len(list); i++ {
		//fmt.Printf("%s\n", list[i])
		result += "        <li>\n"
		parts := strings.Split(list[i], ",")
		// ie. timestamp,shannon,000000,ffffff,ffffff,000000,Them,Us,0,0, 10, 8,  0
		//     0         1       2      3      4      5      6    7  8 9  10  11  12
		// time keeper colorA1 colorA2 colorB1 colorB2 nameA nameB setsA setsB scoreA scoreB possesion
		// 0    1      2       3       4       5       6     7     8     9     10     11     12
		// time keeper colorA1 colorA2 colorB1 colorB2 nameA nameB setsA setsB scoreA scoreB possesion font zoom     sets5    setsShow
		// 0    1      2       3       4       5       6     7     8     9     10     11     12        13   14       15       16
		//                                                                                   1|2       str  zoomOn|  sets5|   setsShowOn|
		//                                                                                                  zoomOff  sets3    setsShowOff
		if len(parts) >= 13 {
			result += parts[0] + ", "
			if parts[12] == "1" {
				result += "*"
			}
			result += parts[6] + ", "
			if parts[12] == "2" {
				result += "*"
			}
			result += parts[7] + ", "

			result += parts[10] + ", "
			result += parts[11] //+ ", "
			// } else if len(parts) == 9 {
			// 	// ie. timestamp,shannon,Them,Us,0,0,10, 8,0
			// 	//     0         1       2    3  4 5 6   7 8
			// 	result += parts[0] + ", "
			// 	if parts[8] == "1" {
			// 		result += "*"
			// 	}
			// 	result += parts[2] + ", "
			// 	if parts[8] == "2" {
			// 		result += "*"
			// 	}
			// 	result += parts[3] + ", "

			// 	result += parts[6] + ", "
			// 	result += parts[7] //+ ", "
		} else {
			// format any links
			line := list[i]
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

func main() {
	VERSION := "2.2"
	MAXLIST := 1000 // max size of list
	// MAXMINUTES := (10)    // max minutes to keep data, per keeper
	// MAXMINUTESALL := (20) // max minutes to keep data, any keeper
	MAXMINUTES := (12 * 60)    // max minutes to keep data, per keeper
	MAXMINUTESALL := (24 * 60) // max minutes to keep data, any keeper
	list := []string{}         // list
	keepers := []Keeper{}
	lastAdd := MAXMINUTESALL // count down from last add

	s := gocron.NewScheduler(time.UTC)
	// check every minute
	s.Cron("*/1 * * * *").Do(func() {
		// fmt.Printf("*/1 * * * *\n")
		// reset by keeper
		for i := 0; i < len(keepers); i++ {
			if keepers[i].lastAdd == 0 {
				fmt.Printf("timed reset for keeper " + keepers[i].name + "\n")
				list = removeKeepersList(keepers[i].name, list)
				keepers[i].lastAdd = -1
			}
			if keepers[i].lastAdd > 0 {
				keepers[i].lastAdd--
			}
		}

		// global reset
		if lastAdd == 0 {
			fmt.Printf("timed reset all\n")
			list = []string{}
			keepers = []Keeper{}
			//lastAdd = -1
		}
		if lastAdd > 0 {
			lastAdd--
		}
	})
	s.StartAsync()

	r := gin.Default()
	r.SetTrustedProxies(nil) // clears error, but not sure if it will block when hosted

	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/style.css", "./resources/style.css")

	// API
	// route: /version
	r.GET("/version", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, VERSION)
	})

	// API
	// route: /reset
	// reset list
	r.GET("/reset", func(c *gin.Context) {
		list = []string{}
		keepers = []Keeper{}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, "")
	})

	// API
	// route: (url)/add?data=(data)
	// add data to list
	r.GET("/add", func(c *gin.Context) {
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
		// if len(parts) == 12 {
		// 	fmt.Printf("long form\n")
		// }
		// if len(parts) == 8 {
		// 	fmt.Printf("long form\n")
		// }
		keeperIndex := getKeepersIndex(parts[0], keepers)
		if keeperIndex != -1 {
			//fmt.Printf("keeper '%s' found\n", parts[0])
			keepers[keeperIndex].lastAdd = MAXMINUTES
		} else {
			k := Keeper{
				name:    parts[0],
				lastAdd: MAXMINUTES,
			}
			keepers = append(keepers, k)
		}
		lastAdd = MAXMINUTESALL

		// prefix time stamp
		currentTime := time.Now()
		entry := currentTime.Format("2006-01-02_15:04:05") + "," + data
		//list = append(list, entry)              // add last
		list = append([]string{entry}, list...) // add first
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, "")

	})

	// DEBUG
	// route: /raw
	// dump all list in raw form
	r.GET("/raw", func(c *gin.Context) {
		msg := strings.Join(list, "\n")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, msg)
	})

	// API
	// route: /json
	// dump all list in raw form
	r.GET("/json", func(c *gin.Context) {
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

		matchList := getKeepersListMany(matchNames, list)
		if count > len(matchList) {
			count = len(matchList)
		}
		partialList := matchList[offset:count]

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.JSON(200, gin.H{"all": partialList})
	})

	// DIRECT USER
	// route: /html
	// html with list of keepers links
	r.GET("/html", func(c *gin.Context) {
		msg := buildKeepersHtml(keepers)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// DIRECT USER
	// route: /
	// html with list of keepers links
	r.GET("/", func(c *gin.Context) {
		msg := buildKeepersHtml(keepers)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// API
	// route: /keepers/json
	// dump keeper names in json form
	r.GET("/keepers/json", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		fmt.Println(keepers)
		names := []string{}
		for i := 0; i < len(keepers); i++ {
			names = append(names, keepers[i].name)
		}
		c.JSON(200, gin.H{"keepers": names})
	})

	// DEBUG
	// route: /(keeper)/raw
	// all raw results for keeper
	r.GET("/:keeperid/raw", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		keeperList := getKeepersList(keeperid, list)
		msg := strings.Join(keeperList, "\n")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, msg)
	})

	// API
	// route: /(keeper)/json
	// json with all results for keeper
	r.GET("/:keeperid/json", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		keeperList := getKeepersList(keeperid, list)

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

	// DIRECT USER
	// route: /(keeper)/html
	// html with results for keeper
	r.GET("/:keeperid/html", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		keeperList := getKeepersList(keeperid, list)
		msg := buildKeeperScoresHtml(keeperid, keeperList)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// DIRECT USER
	// route: /(keeper)/reset
	// reset keeper
	r.GET("/:keeperid/reset", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		list = removeKeepersList(keeperid, list)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, "")
	})

	// DIRECT USER
	// route: /(keeper)
	// html with results for keeper
	r.GET("/:keeperid", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		keeperList := getKeepersList(keeperid, list)
		msg := buildKeeperScoresHtml(keeperid, keeperList)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Data(200, "text/html; charset=utf-8", []byte(msg))
	})

	// API
	// route: /(keeper)/count
	// raw with count of keeper results
	r.GET("/:keeperid/count", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		count := getKeepersCount(keeperid, list)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.String(200, string(rune(count)))
	})

	// API
	// route: /(keeper)/(index)/json
	// json with list[index] results of keeper
	r.GET("/:keeperid/:indexid/json", func(c *gin.Context) {
		keeperid := c.Param("keeperid")
		indexid := c.Param("indexid")
		index, _ := strconv.Atoi(indexid)
		keeperList := getKeepersList(keeperid, list)
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

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	port = ":" + port

	r.Run(port)
}
