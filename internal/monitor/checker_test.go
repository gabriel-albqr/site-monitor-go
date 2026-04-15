package monitor

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"site-monitor-go/internal/domain"
)

func TestChecker_Check_SuccessWithMethodAndHeaders(t *testing.T) {
	t.Parallel()

	var receivedMethod string
	var receivedUserAgent string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		receivedUserAgent = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	checker := NewChecker(5 * time.Second)
	result := checker.Check(domain.Site{
		Name:   "test",
		URL:    srv.URL,
		Method: http.MethodHead,
		Headers: map[string]string{
			"User-Agent": "site-monitor-go-test",
		},
	})

	if result.Status != domain.CheckStatusOnline {
		t.Fatalf("status esperado online, valor: %s", result.Status)
	}
	if result.HTTPStatus != http.StatusNoContent {
		t.Fatalf("HTTP status inesperado: %d", result.HTTPStatus)
	}
	if receivedMethod != http.MethodHead {
		t.Fatalf("método recebido inesperado: %s", receivedMethod)
	}
	if receivedUserAgent != "site-monitor-go-test" {
		t.Fatalf("header User-Agent inesperado: %s", receivedUserAgent)
	}
	if result.ResponseTime <= 0 {
		t.Fatalf("tempo de resposta deveria ser maior que zero")
	}
}

func TestChecker_Check_DefaultMethodIsGET(t *testing.T) {
	t.Parallel()

	var receivedMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewChecker(5 * time.Second)
	result := checker.Check(domain.Site{Name: "test", URL: srv.URL})

	if result.Status != domain.CheckStatusOnline {
		t.Fatalf("status esperado online, valor: %s", result.Status)
	}
	if receivedMethod != http.MethodGet {
		t.Fatalf("método padrão esperado GET, valor: %s", receivedMethod)
	}
}

func TestChecker_Check_UsesSiteTimeout(t *testing.T) {
	t.Parallel()

	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		time.Sleep(1200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewChecker(3 * time.Second)
	result := checker.Check(domain.Site{
		Name:           "slow",
		URL:            srv.URL,
		TimeoutSeconds: 1,
	})

	if result.Status != domain.CheckStatusOffline {
		t.Fatalf("status esperado offline em timeout, valor: %s", result.Status)
	}
	if !strings.Contains(strings.ToLower(result.ErrorMessage), "timeout") {
		t.Fatalf("mensagem de erro deveria indicar timeout, valor: %s", result.ErrorMessage)
	}
	if calls.Load() == 0 {
		t.Fatalf("servidor deveria ter sido chamado")
	}
}

func TestChecker_Check_InvalidURL(t *testing.T) {
	t.Parallel()

	checker := NewChecker(5 * time.Second)
	result := checker.Check(domain.Site{Name: "invalid", URL: "://bad-url"})

	if result.Status != domain.CheckStatusOffline {
		t.Fatalf("status esperado offline, valor: %s", result.Status)
	}
	if !strings.Contains(result.ErrorMessage, "erro ao criar request") {
		t.Fatalf("erro esperado ao criar request, valor: %s", result.ErrorMessage)
	}
}

func TestChecker_Check_ConnectionFailure(t *testing.T) {
	t.Parallel()

	checker := NewChecker(1 * time.Second)
	result := checker.Check(domain.Site{Name: "fail", URL: "http://127.0.0.1:1"})

	if result.Status != domain.CheckStatusOffline {
		t.Fatalf("status esperado offline, valor: %s", result.Status)
	}
	if !strings.Contains(result.ErrorMessage, "falha na request") && !strings.Contains(strings.ToLower(result.ErrorMessage), "timeout") {
		t.Fatalf("erro esperado de falha de request ou timeout, valor: %s", result.ErrorMessage)
	}
}
