package main

import (
	"context"
	"followservice/cmd/provider"
	"followservice/internal/api"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type app struct {
	ctx              context.Context
	cancel           context.CancelFunc
	configuringTasks sync.WaitGroup
	runningTasks     sync.WaitGroup
	env              string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	env := strings.TrimSpace(os.Getenv("ENVIRONMENT"))

	app := &app{
		ctx:    ctx,
		cancel: cancel,
		env:    env,
	}

	app.configuringLog()

	log.Info().Msgf("Starting FollowService service in [%s] enviroment...\n", env)

	provider := provider.NewProvider(env)

	apiEnpoint := provider.ProvideApiEndpoint()

	app.runServerTasks(apiEnpoint)
}

func (app *app) configuringLog() {
	if app.env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = log.With().Caller().Logger()
}

func (app *app) runServerTasks(apiEnpoint *api.Api) {
	app.runningTasks.Add(1)
	go app.runApiEndpoint(apiEnpoint)

	blockForever()

	app.shutdown()
}

func (app *app) runApiEndpoint(apiEnpoint *api.Api) {
	defer app.runningTasks.Done()

	err := apiEnpoint.Run(app.ctx)
	if err != nil {
		log.Panic().Err(err).Msg("Closing FollowService Api failed")
	}
	log.Info().Msg("FollowService Api stopped")
}

func blockForever() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	<-signalCh
}

func (app *app) shutdown() {
	app.cancel()
	log.Info().Msg("Shutting down FollowService Service...")
	app.runningTasks.Wait()
	log.Info().Msg("FollowService Service stopped")
}
