package main

import "github.com/fiufit/metrics/server"

func main() {
	srv := server.NewServer()
	srv.Run()
}
