package main

import (
	// stdlib imports
	"fmt"
	"log"

	// internal imports
	"github.com/gouthamkrishnakv/zerocounter/config"
	"github.com/gouthamkrishnakv/zerocounter/database"
	"github.com/gouthamkrishnakv/zerocounter/logging"
	zc_server "github.com/gouthamkrishnakv/zerocounter/server"
)

// -- constants --

// startupText shows the text to be shown during application statup
const startupText = ` ______                                    _            
|___  /                                   | |           
   / / ___ _ __ ___   ___ ___  _   _ _ __ | |_ ___ _ __ 
  / / / _ \ '__/ _ \ / __/ _ \| | | | '_ \| __/ _ \ '__|
 / /_|  __/ | | (_) | (_| (_) | |_| | | | | ||  __/ |   
/_____\___|_|  \___/ \___\___/ \__,_|_| |_|\__\___|_|   `

// -- functions --

// startup runs the startup method, to set up logging, load configuration
// and set up database
func startup() (error, string) {
	// print startup text
	fmt.Println(startupText)

	// initialize logging
	if loggingErr := logging.Initialize(); loggingErr != nil {
		// DON'T change this. you can't use logging.L if logger initialization
		// failed.
		log.Fatalf("logging.Initialize: %v", loggingErr)
	}

	// initialize configuration
	if configErr := config.Initialize(); configErr != nil {
		return configErr, "config.Initialize failed"
	}

	// initialize database
	if databaseErr := database.Initialize(); databaseErr != nil {
		return databaseErr, "database.Initialize failed"
	}

	return nil, ""
}

// main method starts the application
func main() {
	// run startup
	if startupErr, startupErrMsg := startup(); startupErr != nil {
		logging.L().Fatal().Stack().Err(startupErr).Msg(startupErrMsg)
	}

	// create new server and serve
	server := zc_server.NewServer()
	if serverErr := server.Serve(); serverErr != nil {
		logging.L().Fatal().Stack().Err(serverErr).Msg("server error")
	}

	// INFO: email client or other functionality goes here
}
