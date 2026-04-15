package monitor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"site-monitor-go/internal/domain"
)

// Runner orquestra ciclos contínuos de monitoramento.
type Runner struct {
	service  *Service
	interval time.Duration
}

// NewRunner cria um runner com intervalo entre ciclos.
func NewRunner(service *Service, interval time.Duration) *Runner {
	return &Runner{
		service:  service,
		interval: interval,
	}
}

// Run executa o monitoramento em ciclos contínuos até o contexto ser cancelado.
func (r *Runner) Run(ctx context.Context, sites []domain.Site) error {
	if r.service == nil {
		return errors.New("runner inválido: service não pode ser nil")
	}
	if r.interval <= 0 {
		return errors.New("runner inválido: intervalo deve ser maior que zero")
	}
	if len(sites) == 0 {
		return errors.New("runner inválido: lista de sites está vazia")
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	cycle := 1
	for {
		runCycle(cycle, time.Now().UTC(), r.interval, r.service.CheckSites(sites))
		cycle++

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func runCycle(cycle int, timestamp time.Time, interval time.Duration, results []domain.CheckResult) {
	printCycleHeader(cycle, timestamp, interval)
	PrintResults(results)
}

func printCycleHeader(cycle int, timestamp time.Time, interval time.Duration) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Ciclo %d | %s | proximo em %s\n", cycle, timestamp.Format(time.RFC3339), interval)
	fmt.Println(strings.Repeat("=", 60))
}
