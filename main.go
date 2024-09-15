package main

import "log"

var (
	Version    string
	Commit     string
	CommitDate string
	TreeState  string
)

func main() {
	log.Printf("Version: %s, Commit: %s, Date: %s, State: %s\n", Version, Commit, CommitDate, TreeState)

	server := NewAPIServer(":8080")
	server.Start()
}
