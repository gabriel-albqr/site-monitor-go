package main

import (
	"fmt"
	"log"
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

	service := monitor.NewService(5 * time.Second)
	results := service.CheckSites(cfg.Sites)
	monitor.PrintResults(results)
}
