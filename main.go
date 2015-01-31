package main

import (
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-skelet!"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 8081,
			Usage: "websockets port",
		},
	}

	app.Action = func(c *cli.Context) {
		var port = c.String("port")
		fmt.Printf("listening on port %v\n", port)

		http.HandleFunc("/ws", Handler)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}

	app.Run(os.Args)
}
