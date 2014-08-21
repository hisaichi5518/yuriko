package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "yuriko"
	app.Version = Version
	app.Usage = ""
	app.Author = "hisaichi5518"
	app.Email = "hisaichi5518@gmail.com"
	app.Commands = Commands

	app.Run(os.Args)
}
