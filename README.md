# site-monitor-go

Um monitor de sites feito em Go, com execução contínua, checagem em paralelo e saída clara no terminal.

Este projeto foi desenvolvido durante meus estudos em Golang, com o objetivo de aplicar conceitos importantes da linguagem, como concorrência, organização de código e consumo de APIs HTTP.

---

## Sobre o projeto

O `site-monitor-go` monitora a disponibilidade de múltiplos sites em intervalos definidos.

A cada ciclo, o sistema:

- realiza as requisições HTTP em paralelo
- calcula o tempo de resposta
- exibe um resumo no terminal
- salva os resultados em arquivo para consulta futura

---

## Funcionalidades

- Monitoramento contínuo com intervalo configurável  
- Execução concorrente usando goroutines  
- Configuração por site (método, timeout e headers opcionais)  
- Validação da configuração antes da execução  
- Saída de terminal simples e fácil de entender (em português)  
- Registro dos resultados em arquivo JSONL  

---

## Estrutura do projeto

```text
site-monitor-go/
├── cmd/
│   └── monitor/
│       └── main.go
├── configs/
│   └── sites.json
├── internal/
│   ├── config/
│   ├── domain/
│   ├── monitor/
│   ├── output/
│   └── persistence/
├── data/
│   └── results.jsonl
└── README.md
```

---

## Configuração

Exemplo de `configs/sites.json`:

```json
{
  "check_interval_seconds": 15,
  "sites": [
    {
      "name": "Google",
      "url": "https://www.google.com"
    },
    {
      "name": "GitHub",
      "url": "https://github.com",
      "timeout_seconds": 8
    }
  ]
}
```

- `check_interval_seconds`: tempo entre cada verificação  
- `timeout_seconds`: opcional por site  
- `method` e `headers` também podem ser configurados, se necessário  

---

## Como executar

```bash
go run ./cmd/monitor
```

---

## Exemplo de saída

```text
[INFO] Monitor de sites iniciado
[INFO] 2 sites monitorados | intervalo: 15s

[VERIFICAÇÃO #1] 15/04/2026 15:34:56 | próxima em 15s
-----------------------------------------------------
[RESUMO] ativos: 2 | fora do ar: 0 | tempo médio: 111.9ms

SITE       SITUAÇÃO   HTTP   TEMPO
----------------------------------
Google     ✔ Ativo    200    123ms
GitHub     ✔ Ativo    200    100ms

[SITUAÇÃO] Todos os sites estão funcionando normalmente
```

---

## Persistência

Os resultados são salvos automaticamente em:

```
data/results.jsonl
```

Cada linha representa o resultado de uma checagem, em formato JSON.

Exemplo:

```json
{"site_name":"Google","http_status":200,"response_time_ms":123}
```

---

## Testes

Para rodar os testes:

```bash
go test ./...
```

---

## Observações

O projeto foi estruturado de forma simples, mas organizada, facilitando manutenção e evolução.