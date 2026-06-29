package main

import (
	"os"

	"traffic-cone/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:]))
}
