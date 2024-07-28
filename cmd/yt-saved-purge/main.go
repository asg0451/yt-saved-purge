package main

import "go.coldcutz.net/go-stuff/utils"

func main() {
	ctx, done, log, err := utils.StdSetup()
	if err != nil {
		panic(err)
	}
	defer done()

	log.InfoContext(ctx, "hello")
}
