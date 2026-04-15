package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"site-monitor-go/internal/domain"
)

// JSONLStore persiste resultados em formato JSONL (uma linha JSON por resultado).
type JSONLStore struct {
	filePath string
	mu       sync.Mutex
}

// NewJSONLStore cria um armazenamento de resultados em arquivo JSONL.
func NewJSONLStore(filePath string) (*JSONLStore, error) {
	if filePath == "" {
		return nil, fmt.Errorf("caminho do arquivo de persistência não pode ser vazio")
	}

	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("erro ao preparar diretório de persistência: %w", err)
		}
	}

	return &JSONLStore{filePath: filePath}, nil
}

// SaveCycle salva os resultados de um ciclo no arquivo JSONL.
func (s *JSONLStore) SaveCycle(cycle int, timestamp time.Time, results []domain.CheckResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de persistência: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, result := range results {
		record := cycleRecord{
			Cycle:          cycle,
			CycleTimestamp: timestamp,
			SiteName:       result.Site.Name,
			URL:            result.Site.URL,
			Method:         result.Site.Method,
			HTTPStatus:     result.HTTPStatus,
			IsOnline:       result.Status == domain.CheckStatusOnline,
			ResponseTimeMS: result.ResponseTime.Milliseconds(),
			CheckedAt:      result.CheckedAt,
			ErrorMessage:   result.ErrorMessage,
			TimeoutSeconds: result.Site.TimeoutSeconds,
			RequestHeaders: result.Site.Headers,
			SiteStatus:     string(result.Status),
		}

		if err := encoder.Encode(record); err != nil {
			return fmt.Errorf("erro ao escrever resultado no JSONL: %w", err)
		}
	}

	return nil
}

type cycleRecord struct {
	Cycle          int               `json:"cycle"`
	CycleTimestamp time.Time         `json:"cycle_timestamp"`
	SiteName       string            `json:"site_name"`
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	HTTPStatus     int               `json:"http_status"`
	IsOnline       bool              `json:"is_online"`
	SiteStatus     string            `json:"site_status"`
	ResponseTimeMS int64             `json:"response_time_ms"`
	CheckedAt      time.Time         `json:"checked_at"`
	ErrorMessage   string            `json:"error_message,omitempty"`
	TimeoutSeconds int               `json:"timeout_seconds,omitempty"`
	RequestHeaders map[string]string `json:"request_headers,omitempty"`
}
