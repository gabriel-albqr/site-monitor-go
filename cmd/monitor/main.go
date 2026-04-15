package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"site-monitor-go/internal/config"
	"site-monitor-go/internal/monitor"
	"site-monitor-go/internal/output"
	"site-monitor-go/internal/persistence"
)

func main() {
	const configPath = "configs/sites.json"
	const resultsFilePath = "data/results.jsonl"

	cfg, err := config.LoadSitesConfig(configPath)
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	const requestTimeout = 5 * time.Second

	console := output.NewConsole(os.Stdout)
	console.PrintMonitorStart(cfg.Sites, cfg.CheckInterval(), requestTimeout)

	service := monitor.NewService(requestTimeout)
	runner := monitor.NewRunner(service, cfg.CheckInterval())
	runner.SetReporter(console)

	store, err := persistence.NewJSONLStore(resultsFilePath)
	if err != nil {
		log.Fatalf("erro ao inicializar persistência: %v", err)
	}
	runner.SetResultStore(store)

	if err := runner.Run(ctx, cfg.Sites); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("erro no monitoramento contínuo: %v", err)
	}

	console.PrintMonitorStop()
}
