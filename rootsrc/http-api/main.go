package main

//the http api is the so called MAGMA , which is a http api that requests 
//xml data like login sessionID and everything

import (
	"context"
	"flag"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/metaorior/metafesl/magma/config"
	"github.com/metaorior/metafesl/magma/server"
)

func main() {
	var (
		configFile string
	)
	
	startConfig()
	startLogger()

	//Background means it will always stay on this go-routine
	ctx := context.Background()

	sv, err := server.New()
	if err != nil {
		zerolog.Print("fatal error with http-api setup", err)
	}

	sv.ListenAndServe(
		config.Config.HTTPBind,
		config.Config.HTTPSBind,
		config.Config.CertificatePath,
		config.Config.PrivateKeyPath,
	)

	zerolog.Print("Http api listenting for magma requests")
	<-ctx.Done()
}

func startConfig() {
	// Custom path to configuration file
	flag.StringVar(&configFile, "config", ".env", "Path to configuration file")
	flag.Parse()

	// Override env variables
	gotenv.Load(configFile)

	// Initialize config.* public variables
	config.LoadToMemory()
}

//setup zeroLog so we can call it later
func setupzeroLog() {
	//logs will write with UNIX time
	//TODO browse its vendor package
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}