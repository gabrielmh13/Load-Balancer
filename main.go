package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL   *url.URL
	Proxy *httputil.ReverseProxy
	Alive bool
	mux   sync.RWMutex
}

type LoadBalancer struct {
	Backends []*Backend
	Current  int
	mux      sync.RWMutex
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	lb := LoadBalancer{}

	for _, u := range urls {
		var backend Backend
		url, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}

		backend.URL = url
		backend.Proxy = httputil.NewSingleHostReverseProxy(url)
		backend.Alive = true

		lb.Backends = append(lb.Backends, &backend)
	}

	return &lb
}

func (b *Backend) SetAlive(isAlive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.Alive = isAlive
}

func (lb *LoadBalancer) NextBackend() {
	lb.Current = (lb.Current + 1) % len(lb.Backends)
}

func (lb *LoadBalancer) GetBackend() int {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	defer lb.NextBackend()

	return lb.Current
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backendIndex := lb.GetBackend()

	lb.Backends[backendIndex].Proxy.ServeHTTP(w, r)
	fmt.Printf("Using %s as backend\n", lb.Backends[backendIndex].URL)
}

func isBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	defer conn.Close()
	return true
}

func (lb *LoadBalancer) HealthCheck() {
	for _, b := range lb.Backends {
		status := isBackendAlive(b.URL)
		b.SetAlive(status)
	}
}

func healthCheck(lb *LoadBalancer) {
	t := time.NewTicker(time.Second * 10)
	for v := range t.C {
		log.Printf("Starting health check: " + v.String())
		lb.HealthCheck()
		log.Print("Health check completed\n\n")
	}
}

func main() {
	urls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := NewLoadBalancer(urls)

	go healthCheck(lb)

	if err := http.ListenAndServe(":8080", lb); err != nil {
		log.Fatal(err)
	}

}
