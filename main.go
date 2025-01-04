package main

import (
	"context"
	"database/sql"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/util/cliutil"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	_ "modernc.org/sqlite"
	"os"
	"os/signal"
	"sync"
	"time"
	"tmbgbot/dbstore"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Starting TMBGBot")
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	database, err := NewDatabase(cfg)
	if err != nil {
		panic(err)
	}

	query := dbstore.New(database)

	ticker := time.NewTicker(24 * time.Hour)

	go notificationScheduler(ctx, cfg, query, ticker)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		logger.Info("Waiting for shutdown signal")
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		ticker.Stop()
	}()
	wg.Wait()
}

type Config struct {
	Database string `envconfig:"GOOSE_DBSTRING"`
	UserName string `envconfig:"BSKY_USERNAME"`
	Password string `envconfig:"BSKY_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)

	return &c, err
}

func NewDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite", config.Database)
}

func notificationScheduler(ctx context.Context, cfg *Config, database dbstore.Querier, ticker *time.Ticker) {
	for _ = range ticker.C {
		song, err := database.GetRandomSong(ctx)

		if err == nil {
			// Send notification
			err = sendSongNotification(ctx, cfg, song)
		}
	}
}

func sendSongNotification(ctx context.Context, cfg *Config, song dbstore.SongFormatter) error {
	post := &bsky.FeedPost{
		Text:      song.Format(),
		CreatedAt: time.Now().Local().Format(time.RFC3339),
	}

	xrpcc := &xrpc.Client{Client: cliutil.NewHttpClient(),
		Host: "https://bsky.social",
		Auth: &xrpc.AuthInfo{Handle: cfg.UserName}}

	// Create Auth
	auth, err := atproto.ServerCreateSession(ctx, xrpcc, &atproto.ServerCreateSession_Input{
		Identifier: xrpcc.Auth.Handle,
		Password:   cfg.Password,
	})

	if err != nil {
		return err
	}
	xrpcc.Auth.Did = auth.Did
	xrpcc.Auth.AccessJwt = auth.AccessJwt
	xrpcc.Auth.RefreshJwt = auth.RefreshJwt

	_, err = atproto.RepoCreateRecord(context.TODO(), xrpcc, &atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       xrpcc.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	})

	return err
}
