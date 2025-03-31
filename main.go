package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiKey    string
	ollamaURL string

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "proxy_requests_total",
			Help: "Total number of requests processed by the proxy",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "proxy_request_duration_seconds",
			Help:    "Duration of proxy requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	_ = godotenv.Load()

	apiKey = os.Getenv("OLLAMA_PROXY_KEY")
	ollamaURL = os.Getenv("OLLAMA_URL")
	if apiKey == "" || ollamaURL == "" {
		log.Fatal("Missing required environment variables: OLLAMA_PROXY_KEY and/or OLLAMA_URL")
	}

	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

// Define a custom handler for /metrics that checks for API key
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	providedKey := r.Header.Get("X-API-Key")
	if providedKey != apiKey {
		log.Printf("Unauthorized request to /metrics from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	promhttp.Handler().ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", proxyHandler)
	// Register the custom handler for /metrics
	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Starting Ollama proxy on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	providedKey := r.Header.Get("X-API-Key")
	if providedKey != apiKey {
		log.Printf("Unauthorized request to %s from %s", r.URL.Path, r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		requestCount.WithLabelValues(r.Method, r.URL.Path, "401").Inc()
		return
	}

	url := ollamaURL + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		requestCount.WithLabelValues(r.Method, r.URL.Path, "500").Inc()
		return
	}
	req.Header = cloneHeaders(r.Header)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error forwarding to Ollama: %v", err)
		http.Error(w, "Bad gateway", http.StatusBadGateway)
		requestCount.WithLabelValues(r.Method, r.URL.Path, "502").Inc()
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	// Logging and metrics
	duration := time.Since(start).Seconds()
	log.Printf("%s %s -> %d (%.3fs)", r.Method, r.URL.Path, resp.StatusCode, duration)

	requestCount.WithLabelValues(r.Method, r.URL.Path, http.StatusText(resp.StatusCode)).Inc()
	requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
}

func cloneHeaders(src http.Header) http.Header {
	dest := make(http.Header)
	for k, vv := range src {
		if strings.ToLower(k) == "x-api-key" {
			continue
		}
		for _, v := range vv {
			dest.Add(k, v)
		}
	}
	return dest
}
