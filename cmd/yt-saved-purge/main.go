package main

import "go.coldcutz.net/go-stuff/utils"

type Options struct{}

func main() {
	ctx, done, log, opts, err := utils.StdSetup[Options]()
	if err != nil {
		panic(err)
	}
	defer done()

	log.InfoContext(ctx, "hello", "opts", opts)
}
