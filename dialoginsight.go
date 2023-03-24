package main

import (
	"dialoginsight/osipcollect"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cristalhq/aconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version = "development"
var Built string
var Branch string
var Revision string
var GoVer string

const (
	signalBuffer = 10
)

type Config struct {
	ListenAddr      string        `default:"127.0.0.1:10337" usage:"Local IP and port for the prometheus exporter to listen on. Metrics are exported under http://listen/metrics" flag:"listen" json:"listen"`
	OpensipsMI      string        `default:"http://127.0.0.1:8888/mi" usage:"URL to the mi_http instance for opensips." flag:"opensips_mi" json:"opensips_mi"`
	InsightLabel    string        `default:"insight" usage:"Dialog value starting prefix to indicate it is an insight value (contains labels to process)." flag:"insight_label" json:"insight_label"`
	ExportProfiles  []string      `usage:"List of dialog profiles to export. Comma Separate. Used if export_all is set to false."`
	Timeout         time.Duration `default:"2s" usage:"Timeout duration for opensips API requests."`
	IdleRemove      time.Duration `default:"1m" usage:"If a metric is idle for this long it will be removed from memory."`
	ExportAll       bool          `default:"true" usage:"Whether or not to export all dialog profiles from the instance." flag:"export_all" json:"export_all"`
	EnableProfiling bool          `default:"false" usage:"Enables access to profiling via http://listen/debug/pprof/"`
}

func main() {
	log.Printf("Version: %v Built: %v Branch: %v Revision: %v GoVer: %v", Version, Built, Branch, Revision, GoVer)

	var cfg Config
	cfgLoader := aconfig.LoaderFor(&cfg, aconfig.Config{
		EnvPrefix: "DIALOGINSIGHT",
		FileFlag:  "config",
	})

	if err := cfgLoader.Load(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting with configuration: %+v", cfg)

	osip, err := osipcollect.NewClient(cfg.OpensipsMI, cfg.InsightLabel, cfg.ExportProfiles, cfg.ExportAll, cfg.Timeout, cfg.IdleRemove)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created API instance for OpenSIPs at", cfg.OpensipsMI)

	prometheus.MustRegister(osip)

	log.Println("Listening on", cfg.ListenAddr)

	mux := http.NewServeMux()

	if cfg.EnableProfiling {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/metrics/", promhttp.Handler())

	server := &http.Server{
		Addr:              cfg.ListenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go log.Fatalf("ListenAndServe error: %v", server.ListenAndServe())

	// Wait on Signals
	signals := make(chan os.Signal, signalBuffer)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
