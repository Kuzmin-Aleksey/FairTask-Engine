package app

import (
	"FairTask_Engine/internal/config"
	"FairTask_Engine/internal/domain/service/products"
	cache "FairTask_Engine/internal/infrastructure/persistence/cache/redis"
	"FairTask_Engine/internal/infrastructure/persistence/postgresql"
	"FairTask_Engine/internal/server"
	"FairTask_Engine/pkg/contextx"
	"FairTask_Engine/pkg/logx"
	"FairTask_Engine/pkg/middlewarex"
	"context"
	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lmittmann/tint"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	l := initLogger(cfg.Debug)

	db, err := postgresql.Connect(cfg.DB)
	if err != nil {
		log.Fatal("connect to mysql fail: ", err)
	}

	defer db.Close()

	c, err := cache.Connect(cfg.Redis)
	if err != nil {
		log.Fatal("connect to cache fail: ", err)
	}
	defer c.Close()

	productsRepo := postgresql.NewProductRepo(db)

	productsService := products.NewProductService(productsRepo)

	httpServer := newHttpServer(l, productsService, cfg.Http)

	go func() {
		if cfg.Http.SSLCertPath != "" && cfg.Http.SSLKeyPath != "" {
			log.Println("Starting https server on", cfg.Http.Address)
			if err := httpServer.ListenAndServeTLS(cfg.Http.SSLCertPath, cfg.Http.SSLKeyPath); err != nil {
				log.Fatal("Listen http: ", err)
			}
			return
		}

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Listen http: ", err)
		}
	}()

	sig := <-shutdown
	log.Println("exit by signal: ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("http server shutdown error: ", err)
	}
	cancel()

	os.Exit(0)
}

func newHttpServer(
	l *slog.Logger,
	productsService *products.ProductService,
	cfg *config.HttpConfig,
) *http.Server {
	productsServer := server.NewProductsServer(productsService)

	s := server.NewServer(
		productsServer,
	)

	rtr := mux.NewRouter()
	s.InitRoutes(rtr)

	rtr.Use(
		middlewarex.TraceId,
		middlewarex.Logger,
		middlewarex.RequestLogging(logx.NewSensitiveDataMasker(), 1000),
		middlewarex.ResponseLogging(logx.NewSensitiveDataMasker(), 1000),
		middlewarex.NoCache,
		middlewarex.Recovery,
	)
	if cfg.HandleTimeoutSec > 0 {
		rtr.Use(middlewarex.WithTimeout(time.Duration(cfg.HandleTimeoutSec) * time.Second))
	}

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      rtr,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSec) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSec) * time.Second,
		ErrorLog:     slog.NewLogLogger(l.Handler(), slog.LevelError),
		BaseContext: func(net.Listener) context.Context {
			return contextx.WithLogger(context.Background(), l)
		},
	}
}

func initLogger(debug bool) *slog.Logger {
	if debug {
		return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			AddSource:   false,
			Level:       slog.LevelDebug,
			ReplaceAttr: nil,
			TimeFormat:  time.StampMilli,
			NoColor:     false,
		}))
	}

	rt, err := rotatelogs.New("logs/%Y-%m-%d.log",
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*15),
	)
	if err != nil {
		log.Fatal(err)
	}

	return slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, rt), &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
