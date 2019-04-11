/*
 *     AwesomeBot
 *     Copyright (C) 2019 Asaf Niv
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 */

package main

import "os"

var homeDir, _ = os.UserHomeDir()
var configDir = homeDir + "/.config/"
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
