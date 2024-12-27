package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"github.com/slewiskelly/ock/cmd/ock/internal/subcommands/get"
	ini "github.com/slewiskelly/ock/cmd/ock/internal/subcommands/init"
	"github.com/slewiskelly/ock/cmd/ock/internal/subcommands/list"
	"github.com/slewiskelly/ock/cmd/ock/internal/subcommands/vet"
)

func init() {
	subcommands.Register(&get.Get{}, "")
	subcommands.Register(&ini.Init{}, "")
	subcommands.Register(&list.List{}, "")
	subcommands.Register(&vet.Vet{}, "")

	subcommands.Register(subcommands.HelpCommand(), "")
}

func main() {
	flag.Parse()

	os.Exit(int(subcommands.Execute(context.Background())))
}
