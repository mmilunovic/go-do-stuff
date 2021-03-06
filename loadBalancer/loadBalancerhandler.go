package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (worker *Worker) isAlive() bool {
	timeout := time.Duration(1 * time.Second)

	log.Println("Started Health Check For:", worker.Route)
	_, err := net.DialTimeout("tcp", worker.Route, timeout)
	if err != nil {
		log.Println(worker.Route, "Is Dead")
		log.Println("Health Check Error:", err)
		worker.Alive = false
		return false
	}

	log.Println(worker.Route, "Is Alive")
	worker.Alive = true
	return true
}

func (workerList *WorkerList) init(workerRoutes []string) {
	log.Println("Creating Worker List For Routes:", workerRoutes)

	for _, workerRoute := range workerRoutes {
		var newWorker Worker

		newWorker.Route = workerRoute
		newWorker.Alive = newWorker.isAlive()

		origin, _ := url.Parse("http://" + workerRoute)
		director := func(req *http.Request) {
			req.Header.Add("X-Forwarded-Host", req.Host)
			req.Header.Add("X-Origin-Host", origin.Host)
			req.URL.Scheme = "http"
			req.URL.Host = origin.Host
		}
		newWorker.ReverseProxy = &httputil.ReverseProxy{Director: director}

		log.Println("Worker", newWorker, "Added To Worker List")
		workerList.Workers = append(workerList.Workers, newWorker)
	}

	workerList.Latest = -1
	log.Println("Successfully Created Worker List:", workerList)
}

func (workerList *WorkerList) nextWorker() int {
	return (workerList.Latest + 1) % len(workerList.Workers)
}

func (workerList *WorkerList) loadBalance(w http.ResponseWriter, r *http.Request) {
	if len(workerList.Workers) > 0 {
		workerCount := 0
		for index := workerList.nextWorker(); workerCount < len(workerList.Workers); index = workerList.nextWorker() {
			if workerList.Workers[index].isAlive() {
				log.Println("Routing Request", r.URL, "To", workerList.Workers[index].Route)

				workerList.Workers[index].ReverseProxy.ServeHTTP(w, r)

				workerList.Latest = index
				log.Println("Updated Latest Worker To:", workerList.Latest)

				return
			}
			workerCount++
			workerList.Latest = workerList.nextWorker()
		}
	}
	log.Println("No Workers Available")
	http.Error(w, "No Workers Available", http.StatusServiceUnavailable)
}
