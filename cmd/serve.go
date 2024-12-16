package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"

	"github.com/itsamirhn/linktrkr/internal"
	"github.com/itsamirhn/linktrkr/internal/bot"
	"github.com/itsamirhn/linktrkr/internal/config"
	"github.com/itsamirhn/linktrkr/internal/models"
	"github.com/itsamirhn/linktrkr/pkg"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Telegram Bot serving",
	Run:   serve,
}

func serve(cmd *cobra.Command, _ []string) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Panic(err)
		}
	}()
	loadConfigOrPanic(cmd)

	notifier := pkg.NewNotifier[models.Click]()
	jwtService := pkg.NewJWTService[models.Redirect](config.GlobalConfig.JWT.Secret)

	r := mux.NewRouter()

	b, err := bot.NewBot(
		jwtService,
	)
	if err != nil {
		logrus.WithError(err).Panic("failed to create bot")
	}

	h := internal.NewRedirectHandler(jwtService, notifier)
	r.Handle("/", b.WebhookHandler())
	r.Handle(h.Path(), h).Methods(http.MethodGet)

	s := &http.Server{
		Addr:    ":" + config.GlobalConfig.Server.ListenPort,
		Handler: r,
	}

	var wg conc.WaitGroup

	wg.Go(func() {
		logrus.Info("Starting bot...")
		b.Start()
	})
	wg.Go(func() {
		logrus.Infof("Starting server on: %s", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			logrus.WithError(err).Error("failed to start redirect server")
		}
	})
	wg.Go(notifier.Run)
	wg.Go(func() {
		sub := notifier.Subscribe()
		for msg := range sub {
			b.HandleRedirectMessage(msg)
		}
	})

	wg.Wait()
}
