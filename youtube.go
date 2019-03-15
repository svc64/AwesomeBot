package main

import(
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"os/exec"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func downloadVideo(vidname string) int {
	if strings.HasPrefix(vidname, "http://") ||
		strings.HasPrefix(vidname, "https://")  {
		// Filter vidname to get just the video ID
		s := vidname
		u, err := url.Parse(s)
		handleGeneralError(err)
		rawQ, _ := url.ParseQuery(u.RawQuery)
		if rawQ["v"] != nil {
			fmt.Println("Found YouTube video ID from URL:", rawQ["v"][0])
			vid := rawQ["v"][0]
			returnValue := ytdl(vid)
			return returnValue
		} else {
			fmt.Println("Can't find YouTube video ID from URL")
		}
	} else {
		vid := searchVideoID(vidname)
		fmt.Println("Downloading YouTube video ID: " + vid)
		returnValue := ytdl(vid)
		return returnValue
	}
	return 0
}

// vid = video ID
func ytdl(vid string) int {
	mkCache := exec.Command("mkdir", ".cache")
	dlCmd := exec.Command("youtube-dl", "https://youtu.be/" + vid, "-f 'best[filesize<50M]'", "-x", "--audio-format", "aac", "-o", ".cache/" + vid + ".mp4")
	mkCache.Run()
	mkCache.Wait()
	err := dlCmd.Run()
	dlCmd.Wait()
	if err != nil && strings.ContainsAny(err.Error(), "requested format not available") {
		return 1
	}
	return 0
}

// This API key is taken by decompiling a "spotify downloader" app that actually downloads from youtube
const developerKey = "AIzaSyBzLRQidJbt4BMB8SS7-6c0Nmw-IEfQ_BA"

func searchVideoID(name string) string {
	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(name).
		MaxResults(1)
	response, err := call.Do()
	handleGeneralError(err)

	// Group video results in a separate list.
	videos := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		}
	}

	return getID(videos)
}

func getID(matches map[string]string) string {
	for id := range matches {
		return id
	}
	return ""
}
