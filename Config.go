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
	"os"
	"runtime"
)

var homeDir, _ = os.UserHomeDir()
var configDir = setConfigDir()
var awesomeConfig = configDir + "AwesomeBot"

// checkConfig checks if the configuration folder exists.
func checkConfig() {
	// check if configDir exists
	if !fileExists(configDir) {
		err := os.Mkdir(configDir, 0700)
		checkGeneralError(err)
	}
	// check if awesomeConfig exists
	if !fileExists(awesomeConfig) {
		err := os.Mkdir(awesomeConfig, 0700)
		checkGeneralError(err)
	}
}

// Set the configuration directory based on which OS the user is running
func setConfigDir() string {
	switch runtime.GOOS {
	case "darwin":
		return homeDir + "/Library/"
	case "windows":
		return homeDir // (it should just make an AwesomeBot directory) I have no idea if that's going to work but whatever, Windows isn't worth my time.
	default:
		return homeDir + "/.config/" // Linux, FreeBSD and all the others
	}
}
