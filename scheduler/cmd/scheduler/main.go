package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Rouzip/spider/scheduler"
)

var (
	redisURL = os.Getenv("DS_REDIS")
	bind     = os.Getenv("DS_BIND")

	stdout = log.New(os.Stdout, "[I] ", log.Lmicroseconds|log.Lshortfile)
	stderr = log.New(os.Stderr, "[E] ", log.Lmicroseconds|log.Lshortfile)
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	if err := scheduler.SetRedis(redisURL); err != nil {
		stderr.Println(err)
		return 1
	}

	http.Handle("/", http.RedirectHandler("https://github.com/Rouzip/spider", 302))
	http.HandleFunc("/URL", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getURL(w, r)
		case "POST":
			postURL(w, r)
		default:
			http.Error(w, http.StatusText(404), 404)
		}
	})
	http.HandleFunc("/HTML", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getHTML(w, r)
		case "POST":
			postHTML(w, r)
		default:
			http.Error(w, http.StatusText(404), 404)
		}
	})

	stdout.Println("Listening on", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		stderr.Println(err)
		return 1
	}

	return 0
}
