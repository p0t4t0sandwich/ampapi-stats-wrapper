package main

import (
	webserver "ampapi-stats-wrapper/src"
	"ampapi-stats-wrapper/src/stats"
	"log"
)

func main() {
	settings := stats.NewSettings()
	server := webserver.NewWebServer(settings)
	log.Fatal(server.Run())
}
