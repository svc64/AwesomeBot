/*
 *     AwesomeBot
 *     Copyright (C) 2019 Asaf Niv
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 */

package main

import (
	"flag"
	"fmt"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func downloadVideo(vidname string) (string, bool) { // the bool returns true if we got the video ID, false if not.
	if strings.HasPrefix(vidname, "http://") ||
		strings.HasPrefix(vidname, "https://") {
		// Filter vidname to get just the video ID
		addr := vidname
		u, err := url.Parse(addr)
		checkGeneralError(err)
		rawQ, _ := url.ParseQuery(u.RawQuery)
		if rawQ["v"] != nil {
			vid := rawQ["v"][0]
			fmt.Println("Found YouTube video ID from URL:" + vid)
			ytdl(vid)
			return vid, true
		} else {
			fmt.Println("Can't find YouTube video ID from URL")
		}
	} else {
		vid := searchVideoID(vidname)
		fmt.Println("Downloading YouTube video ID: " + vid)
		ytdl(vid)
		return vid, true
	}
	return "", false
}

// vid = video ID
func ytdl(vid string) {
	homeDir, err := os.UserHomeDir()
	checkGeneralError(err)
	mkCache := exec.Command("mkdir", homeDir+"/.cache")
	dlCmd := exec.Command("youtube-dl", "https://youtu.be/"+vid, "-x", "-f", "best[filesize<800M]", "--audio-format", "aac", "-o", ".cache/"+vid+".mp4")
	err = mkCache.Run()
	checkGeneralError(err)
	err = dlCmd.Run()
	checkGeneralError(err)
}

func searchVideoID(name string) string {
	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: youtubeAPIKey},
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
	checkGeneralError(err)
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

func sendSong(b *tb.Bot, videoID string, m *tb.Message) (*tb.Message, error) {
	filename := ".cache/" + videoID + ".mp4.aac"
	/* Some songs are getting an ".mp4.aac" file extension and some don't
	so we'll check for that and send a .aac file if it exists. */
	if fileExists(filename) {
		file := &tb.Audio{File: tb.FromDisk(filename)}
		m, err := b.Reply(m, file)
		checkError(err, m)
		return m, err
	} else { // song.mp4.aac doesn't exist so we'll try .aac
		filename = ".cache/" + videoID + ".aac"
		file := &tb.Audio{File: tb.FromDisk(filename)}
		m, err := b.Reply(m, file)
		checkError(err, m)
		return m, err
	}
}
