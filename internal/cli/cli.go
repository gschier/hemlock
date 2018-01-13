package cli

import (
	"github.com/alecthomas/kingpin"
	"github.com/tcnksm/go-input"
	"os"
)

var Cmd = kingpin.New("hemlock", "")
var Command = Cmd.Command

var UI = &input.UI{
	Writer: os.Stdout,
	Reader: os.Stdin,
}
