package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
	"github.com/zekrotja/ken/examples/cmdinfo/commands"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	session, err := discordgo.New("")
	if err != nil {
		panic(err)
	}

	k, err := ken.New(session)
	must(err)

	must(k.RegisterCommands(
		new(commands.TestCommand),
	))

	fmt.Println(k.GetCommandInfo().String())
}
