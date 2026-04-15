package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"site-monitor-go/internal/domain"
)

// SitesConfig representa o conteúdo de configuração do arquivo JSON.
type SitesConfig struct {
	CheckIntervalSeconds int           `json:"check_interval_seconds"`
	Sites                []domain.Site `json:"sites"`
}

const defaultHTTPMethod = http.MethodGet

// CheckInterval retorna o intervalo configurado para o ciclo de monitoramento.
func (c SitesConfig) CheckInterval() time.Duration {
	return time.Duration(c.CheckIntervalSeconds) * time.Second
}

// LoadSitesConfig carrega e valida o arquivo de configuração de sites.
func LoadSitesConfig(path string) (SitesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return SitesConfig{}, fmt.Errorf("arquivo de configuração não encontrado: %s", path)
		}

		return SitesConfig{}, fmt.Errorf("erro ao ler arquivo de configuração: %w", err)
	}

	var cfg SitesConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return SitesConfig{}, fmt.Errorf("JSON de configuração inválido: %w", err)
	}

	normalizeSitesConfig(&cfg)

	if err := validateSitesConfig(cfg); err != nil {
		return SitesConfig{}, err
	}

	return cfg, nil
}

func normalizeSitesConfig(cfg *SitesConfig) {
	for i := range cfg.Sites {
		cfg.Sites[i].Name = strings.TrimSpace(cfg.Sites[i].Name)
		cfg.Sites[i].URL = strings.TrimSpace(cfg.Sites[i].URL)

		method := strings.TrimSpace(cfg.Sites[i].Method)
		if method == "" {
			cfg.Sites[i].Method = defaultHTTPMethod
			continue
		}

		cfg.Sites[i].Method = strings.ToUpper(method)
	}
}

func validateSitesConfig(cfg SitesConfig) error {
	if cfg.CheckIntervalSeconds <= 0 {
		return errors.New("configuração inválida: check_interval_seconds deve ser maior que zero")
	}

	if len(cfg.Sites) == 0 {
		return errors.New("configuração inválida: lista de sites está vazia")
	}

	for i, site := range cfg.Sites {
		if strings.TrimSpace(site.Name) == "" {
			return fmt.Errorf("configuração inválida: nome vazio no item %d", i)
		}

		if strings.TrimSpace(site.URL) == "" {
			return fmt.Errorf("configuração inválida: URL vazia no item %d", i)
		}

		if !isValidHTTPMethod(site.Method) {
			return fmt.Errorf("configuração inválida: método HTTP inválido no item %d: %s", i, site.Method)
		}

		if site.TimeoutSeconds < 0 {
			return fmt.Errorf("configuração inválida: timeout_seconds inválido no item %d", i)
		}

		for key := range site.Headers {
			if strings.TrimSpace(key) == "" {
				return fmt.Errorf("configuração inválida: header com chave vazia no item %d", i)
			}
		}
	}

	return nil
}

func isValidHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions:
		return true
	default:
		return false
	}
}
