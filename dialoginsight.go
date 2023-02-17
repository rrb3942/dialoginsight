package main

import (
	"dialoginsight/osipcollect"
	"github.com/cristalhq/aconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	ListenAddr     string        `default:"127.0.0.1:10337" usage:"Local IP and port for the prometheus exporter to listen on." flag:"listen" json:"listen"`
	OpensipsMI     string        `default:"http://127.0.0.1:8888/mi" usage:"url to the mi_http instance for opensips." flag:"opensips_mi" json:"opensips_mi"`
	ExportAll      bool          `default:"true" usage:"Whether or not to export all dialog profiles from the instance." flag:"export_all" json:"export_all"`
	ExportProfiles []string      `usage:"List of Insight dialog profiles to export. Used if export_all is set to false."`
	InsightLabel   string        `default:"insight" usage:"Dialog value starting label to indicate it is an insight value (contains labels to process)." flag:"insight_label" json:"insight_label"`
	Timeout        time.Duration `default:"2s" usage:"Timeout duration for opensips API requests."`
	IdleRemove     time.Duration `default:"1m" usage:"If a metric is idle for this long it will be removed from memory."`
}

func main() {
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

	go log.Fatalf("ListenAndServe error: %v", http.ListenAndServe(cfg.ListenAddr, promhttp.Handler()))

	//Wait on Signals
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
