package main

import "github.com/worldiety/gapt"

func main() {
	cfg := &gapt.Config{}
	cfg.Default()
	cfg.WriteFile("gapt.yml")
}
