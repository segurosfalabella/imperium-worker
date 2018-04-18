package main

import (
	"fmt"
	"imperium-worker/app"
)

func main() {
	err := app.Start("127.0.0.1:8000")

	if err != nil {
		fmt.Println(err.Error())
	}
}
