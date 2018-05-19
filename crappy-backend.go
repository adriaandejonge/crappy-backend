package main

import (
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{ "text": "Response from crappy backend"}`))
	})
	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 10)
		w.Write([]byte(`{ "text": "Response after 10 sec of idle waiting"}`))
	})
	http.HandleFunc("/high-cpu", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan int)
		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
					}
				}
			}()
		}
		time.Sleep(time.Second * 10)
		close(done)
		w.Write([]byte(`{ "text": "Response after 10 sec of high CPU"}`))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
