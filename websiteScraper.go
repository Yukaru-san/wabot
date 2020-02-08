package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetHTMLFromURL returns the whole Body of the given website
func GetHTMLFromURL(url string) string {
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(html)
}

// GetAiringAnime uses proxer.me to find out what will be released
// - day can range from 0 to 1 and stands for today + x days
// (a 0 means: Anime from today)
func GetAiringAnime(day int) string {
	// Starts with index 1
	day++

	// Get html
	html := GetHTMLFromURL("https://proxer.me/calendar#top") // use [1] to [7]
	html = strings.ReplaceAll(html, "\n", "")

	// Categorize into day and anime per day
	dates := strings.Split(html, "<div align=\"left\" class=\"headDate\">")
	animeList := strings.Split(dates[day], "data-ajax=\"true\"") // use [1] to [x]

	var animeInfos []string
	infoDate := ""

	// Every Anime on that day
	for i := 1; i < len(animeList); i++ {
		// name
		nameStartIndex := strings.Index(animeList[i], "\">") + 2
		nameEndIndex := strings.Index(animeList[i][nameStartIndex:], "</a>") + nameStartIndex
		name := animeList[i][nameStartIndex:nameEndIndex]

		// episode
		episodeStartIndex := strings.Index(animeList[i], "Episode")
		episodeEndIndex := strings.Index(animeList[i][episodeStartIndex:], "</a>") + episodeStartIndex
		episode := animeList[i][episodeStartIndex:episodeEndIndex]

		// date & time
		dateStartHelpIndex := strings.Index(animeList[i], "id=\"upTime")
		dateStartIndex := strings.Index(animeList[i][dateStartHelpIndex:], "\">") + 2 + dateStartHelpIndex
		dateEndIndex := strings.Index(animeList[i][dateStartIndex:], "</span>") + dateStartIndex
		unixTime, _ := strconv.ParseInt(animeList[i][dateStartIndex:dateEndIndex], 10, 64)
		dateTime := time.Unix(unixTime, 0).String()

		// single-time-date
		if i == 1 {
			infoDate = fmt.Sprintf("%s/%s/%s", dateTime[8:10], dateTime[5:7], dateTime[0:4])
		}

		// name, episode, day/month/year, time
		//	animeDescription := fmt.Sprintf("%s\nEpisode: %s\n%s/%s/%s\n%s Uhr", name, episode, dateTime[8:10], dateTime[5:7], dateTime[0:4], dateTime[11:16]) writes date aswell
		animeDescription := fmt.Sprintf("%s\n%s\n%s Uhr", name, episode, dateTime[11:16])
		animeInfos = append(animeInfos, animeDescription)
	}

	output := fmt.Sprintf("*Anime-List für den %s:*\n\n\n", infoDate)

	for i := 1; i < len(animeInfos); i++ {
		output += animeInfos[i]

		if i < len(animeInfos)-1 {
			output += "\n\n"
		}
	}

	return output
}
