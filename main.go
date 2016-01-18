package main


import (
  "github.com/elastic/beats/libbeat/beat"
  filesizebeat "./beat"
)

var Version = "0.0.1"
var Name = "filesizebeat"

func main() {
  beat.Run(Name, Version, filesizebeat.New())
}
