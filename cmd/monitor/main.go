package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"site-monitor-go/internal/config"
	"site-monitor-go/internal/monitor"
)

func main() {
	const configPath = "configs/sites.json"

	cfg, err := config.LoadSitesConfig(configPath)
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	fmt.Println("Site Monitor iniciado")
	fmt.Printf("Sites carregados (%d):\n", len(cfg.Sites))

	for i, site := range cfg.Sites {
		fmt.Printf("  %d. %s -> %s\n", i+1, site.Name, site.URL)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	service := monitor.NewService(5 * time.Second)
	runner := monitor.NewRunner(service, cfg.CheckInterval())
	if err := runner.Run(ctx, cfg.Sites); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("erro no monitoramento contínuo: %v", err)
	}

	fmt.Println("Encerrando monitoramento")
}
