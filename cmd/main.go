package main

import (
	"fmt"
	"github.com/worldiety/gapt"
	"os"
)

func main() {
	cfg := &gapt.Config{}
	cfg.Default()
	cfg.WriteFile("gapt.yml")

	fmt.Println(os.FileMode(os.ModeDir).IsRegular())
	//http.Handle("/", http.FileServer(http.Dir("/tmp")))

}
