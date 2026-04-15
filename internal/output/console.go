package output

import (
	"fmt"
	"io"
	"strings"
	"time"

	"site-monitor-go/internal/domain"
)

const separatorLength = 78
const brDateTimeLayout = "02/01/2006 15:04:05"

// Console padroniza a saída do monitor no terminal.
type Console struct {
	out      io.Writer
	location *time.Location

	hasPreviousAverage bool
	previousAverage    time.Duration
	previousBySite     map[string]domain.CheckStatus
	previousTimeBySite map[string]time.Duration
}

// NewConsole cria um escritor de saída formatada.
func NewConsole(out io.Writer) *Console {
	return &Console{
		out:                out,
		location:           recifeLocation(),
		previousBySite:     make(map[string]domain.CheckStatus),
		previousTimeBySite: make(map[string]time.Duration),
	}
}

// PrintMonitorStart exibe o cabeçalho inicial da execução.
func (c *Console) PrintMonitorStart(sites []domain.Site, interval time.Duration, timeout time.Duration) {
	c.printf("[INFO] Site Monitor iniciado\n")
	c.printf("[INFO] Sites configurados: %d | intervalo: %s | tempo limite por site: %s\n", len(sites), interval, timeout)
	for i, site := range sites {
		c.printf("[%02d] %s -> %s\n", i+1, site.Name, site.URL)
	}
	c.printf("\n")
}

// PrintCycleStart exibe informações de início do ciclo.
func (c *Console) PrintCycleStart(cycle int, timestamp time.Time, interval time.Duration) {
	localTimestamp := timestamp.In(c.location)
	formattedTimestamp := localTimestamp.Format(brDateTimeLayout)

	c.printf("\n[VERIFICACAO #%d] %s proxima em %s\n", cycle, formattedTimestamp, interval)
	c.printf("%s\n", strings.Repeat("-", separatorLength))
	c.printf("\n")
}

// PrintCycleResults exibe os resultados consolidados da checagem.
func (c *Console) PrintCycleResults(results []domain.CheckResult) {
	upCount, downCount, average := summarizeResults(results)
	averageDelta := "-"
	if c.hasPreviousAverage {
		averageDelta = formatSignedDuration(average - c.previousAverage)
	}

	c.printf("[RESUMO] ativos: %d | fora do ar: %d | tempo medio: %s | variacao media: %s\n", upCount, downCount, formatDuration(average), averageDelta)
	c.printf("\n")
	c.printf("%-22s %-14s %-6s %-10s %-10s %s\n", "SITE", "SITUACAO", "HTTP", "TEMPO", "VARIACAO", "OBS")
	c.printf("%s\n", strings.Repeat("-", separatorLength))

	for _, result := range results {
		c.printSiteResult(result)
	}
	c.printf("\n")

	if downCount == 0 {
		c.printf("[SITUACAO] Todos os sites estao funcionando normalmente\n")
	} else {
		c.printf("[SITUACAO] Alerta: %d site(s) fora do ar. Verifique os itens marcados com ✖\n", downCount)
	}
	c.printf("\n")

	c.hasPreviousAverage = true
	c.previousAverage = average
	c.storeCurrentStatuses(results)
}

// PrintMonitorStop exibe mensagem de encerramento.
func (c *Console) PrintMonitorStop() {
	c.printf("\n[INFO] Encerrando monitoramento\n")
}

func (c *Console) printSiteResult(result domain.CheckResult) {
	statusLabel := "✖ Fora do ar"
	if result.Status == domain.CheckStatusOnline {
		statusLabel = "✔ Ativo"
	}

	httpStatus := "-"
	if result.HTTPStatus > 0 {
		httpStatus = fmt.Sprintf("%d", result.HTTPStatus)
	}

	errorText := "-"
	if result.ErrorMessage != "" {
		errorText = result.ErrorMessage
	}

	delta := c.deltaForSite(result)
	note := c.transitionNote(result)
	if note == "-" && errorText != "-" {
		note = errorText
	}

	c.printf("%-22s %-14s %-6s %-10s %-10s %s\n", result.Site.Name, statusLabel, httpStatus, formatDuration(result.ResponseTime), delta, note)
}

func formatDuration(duration time.Duration) string {
	return duration.Round(100 * time.Microsecond).String()
}

func formatSignedDuration(duration time.Duration) string {
	rounded := duration.Round(100 * time.Microsecond)
	if rounded > 0 {
		return "+" + rounded.String()
	}

	return rounded.String()
}

func summarizeResults(results []domain.CheckResult) (upCount int, downCount int, average time.Duration) {
	if len(results) == 0 {
		return 0, 0, 0
	}

	var total time.Duration
	for _, result := range results {
		total += result.ResponseTime
		if result.Status == domain.CheckStatusOnline {
			upCount++
			continue
		}

		downCount++
	}

	average = total / time.Duration(len(results))
	return upCount, downCount, average
}

func (c *Console) deltaForSite(result domain.CheckResult) string {
	key := siteKey(result.Site)
	previous, ok := c.previousTimeBySite[key]
	if !ok {
		return "-"
	}

	return formatSignedDuration(result.ResponseTime - previous)
}

func (c *Console) transitionNote(result domain.CheckResult) string {
	key := siteKey(result.Site)
	previous, ok := c.previousBySite[key]
	if !ok {
		return "-"
	}

	if previous == result.Status {
		return "-"
	}

	return shortStatusLabel(previous) + " -> " + shortStatusLabel(result.Status)
}

func (c *Console) storeCurrentStatuses(results []domain.CheckResult) {
	for _, result := range results {
		c.previousBySite[siteKey(result.Site)] = result.Status
		c.previousTimeBySite[siteKey(result.Site)] = result.ResponseTime
	}
}

func siteKey(site domain.Site) string {
	return site.Name + "|" + site.URL
}

func shortStatusLabel(status domain.CheckStatus) string {
	if status == domain.CheckStatusOnline {
		return "Ativo"
	}

	return "Fora do ar"
}

func recifeLocation() *time.Location {
	location, err := time.LoadLocation("America/Recife")
	if err == nil {
		return location
	}

	return time.FixedZone("America/Recife", -3*60*60)
}

func (c *Console) printf(format string, args ...any) {
	_, _ = fmt.Fprintf(c.out, format, args...)
}
