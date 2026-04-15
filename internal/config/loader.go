package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"site-monitor-go/internal/domain"
)

// SitesConfig representa o conteúdo de configuração do arquivo JSON.
type SitesConfig struct {
	Sites []domain.Site `json:"sites"`
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

	if err := validateSitesConfig(cfg); err != nil {
		return SitesConfig{}, err
	}

	return cfg, nil
}

func validateSitesConfig(cfg SitesConfig) error {
	if len(cfg.Sites) == 0 {
		return errors.New("configuração inválida: lista de sites está vazia")
	}

	for i, site := range cfg.Sites {
		if strings.TrimSpace(site.URL) == "" {
			return fmt.Errorf("configuração inválida: URL vazia no item %d", i)
		}
	}

	return nil
}
