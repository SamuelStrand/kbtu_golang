package exchange

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetRate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/convert", r.URL.Path)
		require.Equal(t, "USD", r.URL.Query().Get("from"))
		require.Equal(t, "EUR", r.URL.Query().Get("to"))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.91}`)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	rate, err := svc.GetRate("USD", "EUR")
	require.NoError(t, err)
	require.InDelta(t, 0.91, rate, 0.0001)
}

func TestGetRate_APIBusinessError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"error":"invalid currency pair"}`)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("AAA", "BBB")
	require.Error(t, err)
	require.Contains(t, err.Error(), "api error: invalid currency pair")
}

func TestGetRate_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":`) // broken
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode error")
}

func TestGetRate_NonJSONBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error: truncated text")
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode error")
}

func TestGetRate_EmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode error")
}

func TestGetRate_Server500UnexpectedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"panic"}`)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected status: 500")
}

func TestGetRate_SlowResponseTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"base":"USD","target":"EUR","rate":0.5}`)
	}))
	defer srv.Close()

	svc := NewExchangeService(srv.URL)
	svc.Client.Timeout = 50 * time.Millisecond // force timeout

	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	require.Contains(t, err.Error(), "network error")
}
