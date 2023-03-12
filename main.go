package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	Proxy []*httputil.ReverseProxy
	URLS  []*url.URL
	index int
	mux   sync.RWMutex
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	lb := LoadBalancer{}

	for _, u := range urls {
		url, err := url.Parse(u)
		if err != nil {
			log.Fatal(err)
		}
		lb.URLS = append(lb.URLS, url)
		lb.Proxy = append(lb.Proxy, httputil.NewSingleHostReverseProxy(url))
	}

	return &lb
}

func (lb *LoadBalancer) NextIndex() {
	lb.index = (lb.index + 1) % len(lb.URLS)
}

func (lb *LoadBalancer) GetIndex() int {
	lb.mux.Lock()
	defer lb.mux.Unlock()

	defer lb.NextIndex()

	return lb.index
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backendIndex := lb.GetIndex()

	lb.Proxy[backendIndex].ServeHTTP(w, r)
	fmt.Printf("Using %s as backend\n", lb.URLS[backendIndex])
}

func main() {
	urls := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := NewLoadBalancer(urls)

	if err := http.ListenAndServe(":8080", lb); err != nil {
		log.Fatal(err)
	}

}
