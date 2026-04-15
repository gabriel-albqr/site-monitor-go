package monitor

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"site-monitor-go/internal/domain"
)

const defaultRequestTimeout = 5 * time.Second

// Checker executa checagens HTTP individuais de forma isolada.
type Checker struct {
	client         *http.Client
	defaultTimeout time.Duration
}

// NewChecker cria um checker com timeout configurado.
func NewChecker(timeout time.Duration) *Checker {
	if timeout <= 0 {
		timeout = defaultRequestTimeout
	}

	return &Checker{
		client:         &http.Client{},
		defaultTimeout: timeout,
	}
}

// Check executa uma checagem HTTP para um site monitorado.
func (c *Checker) Check(site domain.Site) domain.CheckResult {
	checkedAt := time.Now().UTC()
	startedAt := time.Now()

	result := domain.CheckResult{
		Site:      site,
		Status:    domain.CheckStatusOffline,
		CheckedAt: checkedAt,
	}

	effectiveTimeout := c.defaultTimeout
	if site.TimeoutSeconds > 0 {
		effectiveTimeout = time.Duration(site.TimeoutSeconds) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), effectiveTimeout)
	defer cancel()

	method := site.Method
	if method == "" {
		method = http.MethodGet
	}

	request, err := http.NewRequestWithContext(ctx, method, site.URL, nil)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("erro ao criar request: %v", err)
		return result
	}

	for key, value := range site.Headers {
		request.Header.Set(key, value)
	}

	response, err := c.client.Do(request)
	result.ResponseTime = time.Since(startedAt)
	if err != nil {
		result.ErrorMessage = formatRequestError(err)
		return result
	}
	defer response.Body.Close()

	result.HTTPStatus = response.StatusCode
	result.Status = domain.CheckStatusOnline
	return result
}

func formatRequestError(err error) string {
	if err == nil {
		return ""
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Sprintf("timeout ao chamar site: %v", err)
	}

	if strings.Contains(strings.ToLower(err.Error()), "timeout") {
		return fmt.Sprintf("timeout ao chamar site: %v", err)
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fmt.Sprintf("timeout ao chamar site: %v", err)
	}

	return fmt.Sprintf("falha na request: %v", err)
}
