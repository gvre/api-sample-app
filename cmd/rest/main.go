package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"

	"github.com/gvre/api-sample-app/cmd/rest/api"
	"github.com/gvre/api-sample-app/user"
)

func main() {
	// Flags
	var (
		host = flag.String("host", "", "listen host")
		port = flag.String("port", "8080", "listen port")
	)
	flag.Parse()

	// skip errors if .env file is not found
	_ = godotenv.Load()

	// Database
	db, err := pgxpool.Connect(context.Background(), "postgres://")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Services
	userService := user.NewService(user.NewDatabaseRepository(db))

	// Logger
	handlerOptions := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.String("time", time.Now().UTC().Format("2006-01-02T15:04:05.000Z"))
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(os.Stdout, &handlerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Rest server
	server := api.NewServer(userService, logger)

	// HTTP server
	srv := &http.Server{
		Addr:    net.JoinHostPort(*host, *port),
		Handler: server.Router,

		// ReadHeaderTimeout is the amount of time allowed to read
		// request headers. The connection's read deadline is reset
		// after reading the headers and the Handler can decide what
		// is considered too slow for the body. If ReadHeaderTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		ReadHeaderTimeout: 5 * time.Second,

		// ReadTimeout is the maximum duration for reading the entire
		// request, including the body.
		//
		// Because ReadTimeout does not let Handlers make per-request
		// decisions on each request body's acceptable deadline or
		// upload rate, most users will prefer to use
		// ReadHeaderTimeout. It is valid to use them both.
		ReadTimeout: 5 * time.Second,

		// WriteTimeout is the maximum duration before timing out
		// writes of the response. It is reset whenever a new
		// request's header is read. Like ReadTimeout, it does not
		// let Handlers make decisions on a per-request basis.
		WriteTimeout: 5 * time.Second,

		// IdleTimeout is the maximum amount of time to wait for the
		// next request when keep-alives are enabled. If IdleTimeout
		// is zero, the value of ReadTimeout is used. If both are
		// zero, there is no timeout.
		IdleTimeout: 5 * time.Second,
	}

	start(srv, *host, *port)
}

func start(srv *http.Server, host, port string) {
	// Shutdown the http server when a signal INT, TERM or QUIT is received.
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(
			signals,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		signal.Ignore(syscall.SIGHUP)

		<-signals
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Println("HTTP server Shutdown:", err)
		}
	}()

	fmt.Println("[API] Started. Address:", net.JoinHostPort(host, port))
	log.Fatal(srv.ListenAndServe())
}
