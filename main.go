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

	ctx = WithLogger(ctx, logger)

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
	logger := FromContext(ctx)
	for _ = range ticker.C {
		song, err := database.GetRandomSong(ctx)

		if err != nil {
			logger.Error("Error Loading Song")
			continue
		}
		// Send notification
		logger.Info("Loaded song", "song", song.Format())
		err = sendSongNotification(ctx, cfg, song)
		if err != nil {
			logger.Error("Error Posting To Bluesky", "error", err)
		}
	}
}

func sendSongNotification(ctx context.Context, cfg *Config, song dbstore.SongFormatter) error {
	logger := FromContext(ctx)

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
		logger.Error("Error Authenticating To Bluesky", "error", err)
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

func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value("logger").(*slog.Logger); ok {
		return logger
	}
	// Hopefully this never happens
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
