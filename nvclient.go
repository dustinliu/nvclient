package main

import (
	"os"

	"github.com/dustinliu/nvclient/pkg/app"
	"github.com/urfave/cli/v2"
)

func main() {
	app := app.NewApp()

	err := app.Run(os.Args)
	if err != nil {
		cli.HandleExitCoder(err)
	}
}
