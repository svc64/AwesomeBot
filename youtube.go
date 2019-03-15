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

func downloadVideo(vidname string) {
	if strings.HasPrefix(vidname, "http") {
		// Filter vidname to get just the video ID
		s := vidname
		u, err := url.Parse(s)
		handleGeneralError(err)
		rawQ, _ := url.ParseQuery(u.RawQuery)
		if rawQ["v"] != nil {
			fmt.Println("Found YouTube video ID from URL:", rawQ["v"][0])
			vid := rawQ["v"][0]
			mkCache := exec.Command("mkdir", ".cache")
			dlCmd := exec.Command("youtube-dl", "https://youtu.be/" + vid, "-x", "--audio-format", "aac", "-o", ".cache/" + vid + ".mp4")
			mkCache.Run()
			mkCache.Wait()
			dlCmd.Run()
			dlCmd.Wait()
		} else {
			fmt.Println("Can't find YouTube video ID from URL")
		}
	} else {
		vid := searchVideoID(vidname)
		fmt.Println("Downloading YouTube video ID: " + vid)
		mkCache := exec.Command("mkdir", ".cache")
		dlCmd := exec.Command("youtube-dl", "https://youtu.be/" + vid, "-x", "--audio-format", "aac", "-o", ".cache/" + vid + ".mp4")
		mkCache.Run()
		mkCache.Wait()
		dlCmd.Run()
		dlCmd.Wait()
	}
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
