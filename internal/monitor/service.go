package monitor

import (
	"fmt"
	"strings"
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

// PrintResults exibe os resultados em formato amigável no terminal.
func PrintResults(results []domain.CheckResult) {
	fmt.Printf("\nResultados das checagens (%d):\n", len(results))
	fmt.Println(strings.Repeat("-", 60))

	for i, result := range results {
		printResult(i+1, result)
		if i < len(results)-1 {
			fmt.Println(strings.Repeat("-", 60))
		}
	}
}

func printResult(index int, result domain.CheckResult) {
	fmt.Printf("%d. %s\n", index, result.Site.Name)
	fmt.Printf("   URL: %s\n", result.Site.URL)
	fmt.Printf("   Status: %s\n", result.Status)
	fmt.Printf("   HTTP: %d\n", result.HTTPStatus)
	fmt.Printf("   Resposta: %s\n", result.ResponseTime)
	fmt.Printf("   Checado em: %s\n", result.CheckedAt.Format(time.RFC3339))
	if result.ErrorMessage != "" {
		fmt.Printf("   Erro: %s\n", result.ErrorMessage)
	}
}
