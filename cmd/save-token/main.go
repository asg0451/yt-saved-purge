package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"go.coldcutz.net/go-stuff/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Options struct {
	ClientSecretJSONPath string `long:"client-secret-json-path" env:"CLIENT_SECRET_JSON_PATH" description:"Path to client secret JSON file" default:"client_secret.json"`
	TokenFilePath        string `long:"token-file-path" env:"TOKEN_FILE_PATH" description:"Path to token file" default:"token.json"`
}

func main() {
	ctx, done, log, opts, err := utils.StdSetup[Options]()
	if err != nil {
		panic(err)
	}
	defer done()

	if err := run(ctx, log, opts); err != nil {
		log.Error("error", "err", err)
	}
}

func run(ctx context.Context, log *slog.Logger, opts Options) error {
	b, err := os.ReadFile(opts.ClientSecretJSONPath)
	if err != nil {
		return fmt.Errorf("reading client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		return fmt.Errorf("parsing client secret file to config: %v", err)
	}
	tok, err := getToken(ctx, config)
	if err != nil {
		return fmt.Errorf("getting token: %v", err)
	}

	if err := saveToken(opts.TokenFilePath, tok); err != nil {
		return fmt.Errorf("saving token: %v", err)
	}

	log.Info("Token saved to file", "path", opts.TokenFilePath)

	return nil
}

func getToken(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("reading authorization code: %v", err)
	}

	tok, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("retrieving token from web: %v", err)
	}
	return tok, nil
}

func saveToken(file string, token *oauth2.Token) error {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("encoding token: %v", err)
	}
	return nil
}
