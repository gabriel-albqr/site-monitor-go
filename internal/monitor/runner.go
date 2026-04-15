package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"site-monitor-go/internal/domain"
)

// CycleReporter define a saída do runner sem acoplar a lógica de monitoramento ao terminal.
type CycleReporter interface {
	PrintCycleStart(cycle int, timestamp time.Time, interval time.Duration)
	PrintCycleResults(results []domain.CheckResult)
}

// CycleResultStore define a persistência dos resultados sem acoplamento ao runner.
type CycleResultStore interface {
	SaveCycle(cycle int, timestamp time.Time, results []domain.CheckResult) error
}

// Runner orquestra ciclos contínuos de monitoramento.
type Runner struct {
	service  *Service
	interval time.Duration
	reporter CycleReporter
	store    CycleResultStore
}

// NewRunner cria um runner com intervalo entre ciclos.
func NewRunner(service *Service, interval time.Duration) *Runner {
	return &Runner{
		service:  service,
		interval: interval,
	}
}

// SetReporter define o reporter de saída do runner.
func (r *Runner) SetReporter(reporter CycleReporter) {
	r.reporter = reporter
}

// SetResultStore define o destino de persistência dos resultados.
func (r *Runner) SetResultStore(store CycleResultStore) {
	r.store = store
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
	if r.reporter == nil {
		return errors.New("runner inválido: reporter não pode ser nil")
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	cycle := 1
	for {
		timestamp := time.Now().UTC()
		results := r.service.CheckSites(sites)

		r.reporter.PrintCycleStart(cycle, timestamp, r.interval)
		r.reporter.PrintCycleResults(results)

		if r.store != nil {
			if err := r.store.SaveCycle(cycle, timestamp, results); err != nil {
				return fmt.Errorf("falha ao persistir resultados do ciclo %d: %w", cycle, err)
			}
		}

		cycle++

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
