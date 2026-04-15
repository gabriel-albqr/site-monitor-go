package monitor

import (
	"sync"
	"time"

	"site-monitor-go/internal/domain"
)

// Service coordena a execução das checagens.
type Service struct {
	checker *Checker
}

// NewService cria um serviço de monitoramento com timeout configurado.
func NewService(timeout time.Duration) *Service {
	return &Service{checker: NewChecker(timeout)}
}

// CheckSites executa as checagens em paralelo com coleta segura.
func (s *Service) CheckSites(sites []domain.Site) []domain.CheckResult {
	if len(sites) == 0 {
		return []domain.CheckResult{}
	}

	type indexedResult struct {
		index  int
		result domain.CheckResult
	}

	results := make([]domain.CheckResult, len(sites))
	resultCh := make(chan indexedResult, len(sites))

	var wg sync.WaitGroup
	for i, site := range sites {
		wg.Add(1)

		go func(index int, currentSite domain.Site) {
			defer wg.Done()
			resultCh <- indexedResult{
				index:  index,
				result: s.checker.Check(currentSite),
			}
		}(i, site)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for item := range resultCh {
		results[item.index] = item.result
	}

	return results
}
