package main

import (
	"os"

	"github.com/dustinliu/nvclient/pkg/nvc"
	"github.com/urfave/cli/v2"
)

func main() {
	app := nvc.NewApp()

	err := app.Run(os.Args)
	if err != nil {
		cli.HandleExitCoder(err)
	}
}
