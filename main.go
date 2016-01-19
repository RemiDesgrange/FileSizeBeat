package main

import (
	filesizebeat "./beat"
	"github.com/elastic/beats/libbeat/beat"
)

var Version = "0.0.1"
var Name = "filesizebeat"

func main() {
	beat.Run(Name, Version, filesizebeat.New())
}
