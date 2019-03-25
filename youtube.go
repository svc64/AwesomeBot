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
	tb "git.asafniv.me/blzit420/telebot.v2"
)

func downloadVideo(vidname string) (string, bool) { // the bool returns true if we got the video ID, false if not.
	var ID string
	if strings.HasPrefix(vidname, "http://") ||
		strings.HasPrefix(vidname, "https://")  {
		// Filter vidname to get just the video ID
		addr := vidname
		u, err := url.Parse(addr)
		handleGeneralError(err)
		rawQ, _ := url.ParseQuery(u.RawQuery)
		if rawQ["v"] != nil {
			ID = rawQ["v"][0]
			fmt.Println("Found YouTube video ID from URL:" + ID)
		} else {
			fmt.Println("Can't find YouTube video ID from URL")
		}
	} else {
		ID = searchVideoID(vidname)
		fmt.Println("Downloading YouTube video ID: " + ID)
	}
	err := ytdl(ID)
	if err != nil {
		return "", false
	} else {
		return ID, true
	}
	return "", false
}

// vid = video ID
func ytdl(vid string) error {
	mkCache := exec.Command("mkdir", ".cache")
	dlCmd := exec.Command("youtube-dl", "https://youtu.be/" + vid, "-x", "-f", "bestaudio[filesize<800M]", "--audio-format", "aac", "-o", ".cache/" + vid + ".mp4")
	err := mkCache.Run()
	handleGeneralError(err)
	if err != nil {
		return err
	}
	err = dlCmd.Run()
	handleGeneralError(err)
	return err
}

// Put your own API key here - this one is restricted to my IP
// If you want a key that doesn't identify you, you can probably get one by decompiling apps and copying the key.
const developerKey = "AIzaSyAcJAjE4cvQsuHfJmE-EVTK5C88DUiPHIg"

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

func sendSong(b *tb.Bot, videoID string, m *tb.Message) {
	filename := ".cache/" + videoID + ".mp4.aac"
	/* Some songs are getting an ".mp4.aac" file extension and some don't
	so we'll check for that and send a .aac file if it exists. */
	if fileExists(filename) {
		file := &tb.Audio{File: tb.FromDisk(filename)}
		_ ,err := b.Reply(m, file)
		handleError(nil, err, *m)
	} else { // song.mp4.aac doesn't exist so we'll try .aac
		filename = ".cache/" + videoID + ".aac"
		file := &tb.Audio{File: tb.FromDisk(filename)}
		_ ,err := b.Reply(m, file)
		handleError(nil, err, *m)
	}
}