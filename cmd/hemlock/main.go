package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/gschier/hemlock/internal/cli"
	"os"

	// Commands
	_ "github.com/gschier/hemlock/internal/cli"
)

var (
	Cmd = kingpin.New("hemlock", "")
)

func main() {
	_, err := cli.Cmd.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
}
