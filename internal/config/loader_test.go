package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSitesConfig_SuccessWithDefaultsAndNormalization(t *testing.T) {
	t.Parallel()

	content := `{
  "check_interval_seconds": 15,
  "sites": [
    {
      "name": "  Google  ",
      "url": "  https://www.google.com  "
    },
    {
      "name": "GitHub",
      "url": "https://github.com",
      "method": "post",
      "timeout_seconds": 8,
      "headers": {
        "User-Agent": "site-monitor-go/1.0"
      }
    }
  ]
}`

	configPath := writeTempConfigFile(t, content)

	cfg, err := LoadSitesConfig(configPath)
	if err != nil {
		t.Fatalf("esperava configuração válida, mas retornou erro: %v", err)
	}

	if cfg.CheckIntervalSeconds != 15 {
		t.Fatalf("intervalo inesperado: %d", cfg.CheckIntervalSeconds)
	}

	if got := len(cfg.Sites); got != 2 {
		t.Fatalf("quantidade de sites inesperada: %d", got)
	}

	if cfg.Sites[0].Name != "Google" {
		t.Fatalf("nome não normalizado, valor: %q", cfg.Sites[0].Name)
	}
	if cfg.Sites[0].URL != "https://www.google.com" {
		t.Fatalf("URL não normalizada, valor: %q", cfg.Sites[0].URL)
	}
	if cfg.Sites[0].Method != "GET" {
		t.Fatalf("método padrão esperado GET, valor: %q", cfg.Sites[0].Method)
	}

	if cfg.Sites[1].Method != "POST" {
		t.Fatalf("método deveria ser normalizado para POST, valor: %q", cfg.Sites[1].Method)
	}
	if cfg.Sites[1].TimeoutSeconds != 8 {
		t.Fatalf("timeout esperado 8, valor: %d", cfg.Sites[1].TimeoutSeconds)
	}
	if cfg.Sites[1].Headers["User-Agent"] != "site-monitor-go/1.0" {
		t.Fatalf("header esperado não encontrado")
	}
}

func TestLoadSitesConfig_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		content    string
		errorMatch string
	}{
		{
			name:       "json inválido",
			content:    `{`,
			errorMatch: "JSON de configuração inválido",
		},
		{
			name: "intervalo inválido",
			content: `{
  "check_interval_seconds": 0,
  "sites": [{"name":"Google","url":"https://www.google.com"}]
}`,
			errorMatch: "check_interval_seconds deve ser maior que zero",
		},
		{
			name: "lista vazia",
			content: `{
  "check_interval_seconds": 10,
  "sites": []
}`,
			errorMatch: "lista de sites está vazia",
		},
		{
			name: "nome vazio",
			content: `{
  "check_interval_seconds": 10,
  "sites": [{"name":" ","url":"https://www.google.com"}]
}`,
			errorMatch: "nome vazio",
		},
		{
			name: "url vazia",
			content: `{
  "check_interval_seconds": 10,
  "sites": [{"name":"Google","url":""}]
}`,
			errorMatch: "URL vazia",
		},
		{
			name: "método inválido",
			content: `{
  "check_interval_seconds": 10,
  "sites": [{"name":"Google","url":"https://www.google.com","method":"FETCH"}]
}`,
			errorMatch: "método HTTP inválido",
		},
		{
			name: "timeout inválido",
			content: `{
  "check_interval_seconds": 10,
  "sites": [{"name":"Google","url":"https://www.google.com","timeout_seconds":-1}]
}`,
			errorMatch: "timeout_seconds inválido",
		},
		{
			name: "header com chave vazia",
			content: `{
  "check_interval_seconds": 10,
  "sites": [{"name":"Google","url":"https://www.google.com","headers":{"":"abc"}}]
}`,
			errorMatch: "header com chave vazia",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			configPath := writeTempConfigFile(t, tt.content)

			_, err := LoadSitesConfig(configPath)
			if err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}
			if !strings.Contains(err.Error(), tt.errorMatch) {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

func TestLoadSitesConfig_FileNotFound(t *testing.T) {
	t.Parallel()

	_, err := LoadSitesConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatalf("esperava erro de arquivo inexistente")
	}
	if !strings.Contains(err.Error(), "arquivo de configuração não encontrado") {
		t.Fatalf("erro inesperado: %v", err)
	}
}

func writeTempConfigFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "sites.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("erro ao criar arquivo temporário: %v", err)
	}

	return path
}
