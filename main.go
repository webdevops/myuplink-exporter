package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	flags "github.com/jessevdk/go-flags"
	cache "github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/webdevops/myuplink-exporter/config"
	"github.com/webdevops/myuplink-exporter/myuplink"
)

const (
	Author    = "webdevops.io"
	UserAgent = "myuplink-exporter/"
)

var (
	argparser *flags.Parser
	opts      config.Opts

	myuplinkClient *myuplink.Client

	globalCache *cache.Cache

	// Git version information
	gitCommit = "<unknown>"
	gitTag    = "<unknown>"
)

func main() {
	initArgparser()
	initLogger()

	logger.Infof("starting myuplink-plug-exporter v%s (%s; %s; by %v)", gitTag, gitCommit, runtime.Version(), Author)
	logger.Info(string(opts.GetJson()))

	globalCache = cache.New(60*time.Minute, 1*time.Minute)
	totalParamCache.Init()

	logger.Infof("connecting to myUplink")
	myuplinkClient = myuplink.NewClient(logger)
	myuplinkClient.SetDebugMode(opts.Logger.Development)
	myuplinkClient.SetApiUrl(opts.MyUplink.Url)
	myuplinkClient.SetUserAgent(UserAgent + gitTag)
	myuplinkClient.SetAuth(opts.MyUplink.Auth.ClientID, opts.MyUplink.Auth.ClientSecret)
	if err := myuplinkClient.Connect(context.Background()); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("starting http server on %s", opts.Server.Bind)
	startHttpServer()
}

// init argparser and parse/validate arguments
func initArgparser() {
	argparser = flags.NewParser(&opts, flags.Default)
	_, err := argparser.Parse()

	// check if there is an parse error
	if err != nil {
		var flagsErr *flags.Error
		if ok := errors.As(err, &flagsErr); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

// start and handle prometheus handler
func startHttpServer() {
	mux := http.NewServeMux()

	// healthz
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "Ok"); err != nil {
			logger.Error(err)
		}
	})

	// readyz
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "Ok"); err != nil {
			logger.Error(err)
		}
	})

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/probe", myuplinkProbe)

	srv := &http.Server{
		Addr:         opts.Server.Bind,
		Handler:      mux,
		ReadTimeout:  opts.Server.ReadTimeout,
		WriteTimeout: opts.Server.WriteTimeout,
	}
	logger.Fatal(srv.ListenAndServe())
}
