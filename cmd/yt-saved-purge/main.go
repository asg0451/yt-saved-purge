package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/davecgh/go-spew/spew"
	"go.coldcutz.net/go-stuff/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Options struct {
	TokenJSON            string `long:"token-json" env:"TOKEN_JSON" description:"Token JSON" required:"true"`
	ClientSecretJSONPath string `long:"client-secret-json-path" env:"CLIENT_SECRET_JSON_PATH" description:"Path to client secret JSON file" default:"client_secret.json"`
}

func main() {
	ctx, done, log, opts, err := utils.StdSetup[Options]()
	if err != nil {
		panic(err)
	}
	defer done()

	if err := run(ctx, log, opts); err != nil {
		log.Error("error", "err", err)
		done()
		os.Exit(1)
	}
}

func run(ctx context.Context, log *slog.Logger, opts Options) error {
	b, err := os.ReadFile(opts.ClientSecretJSONPath)
	if err != nil {
		return fmt.Errorf("reading client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		return fmt.Errorf("parsing client secret file to config: %v", err)
	}

	t := &oauth2.Token{}
	if err := json.Unmarshal([]byte(opts.TokenJSON), t); err != nil {
		return fmt.Errorf("unmarshalling token: %v", err)
	}

	client := config.Client(ctx, t)

	service, err := youtube.New(client)
	if err != nil {
		return fmt.Errorf("creating youtube service: %v", err)
	}

	res, err := service.Playlists.List([]string{"id", "snippet"}).Mine(true).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("getting playlists: %v", err)
	}
	spew.Dump(res)

	return nil
}
