package main

import (
        "os"

        _ "github.com/dapicard/filebeat.mongodb.output/mongodb"

        "github.com/elastic/beats/filebeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}