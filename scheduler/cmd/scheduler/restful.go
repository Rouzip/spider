package main

import (
	"bufio"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/Rouzip/spider/scheduler"
)

var (
	urlWait  = make(chan struct{})
	htmlWait = make(chan struct{})

	urlFb, htmlFb     bool
	urlJobs, htmlJobs uint
	urlMu, htmlMu     sync.RWMutex
)

func isOwari() bool {
	urlMu.RLock()
	htmlMu.RLock()

	owari := urlFb && htmlFb &&
		urlJobs == 0 && htmlJobs == 0

	htmlMu.RUnlock()
	urlMu.RUnlock()

	return owari
}

func getURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

REFETCH:
	url, isDrained, err := scheduler.PopURL(r.Context())
	if isDrained {
		if isOwari() {
			stdout.Println("OWATTA! (from getURL)")
			w.WriteHeader(410)
			return
		}
		stdout.Println("URL temporarily drained, waiting.")
		select {
		case <-urlWait:
			goto REFETCH
		case <-r.Context().Done():
			return
		}
	}
	if err != nil {
		stderr.Println("getURL:", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(url)))
	w.WriteHeader(200)
	io.WriteString(w, url)

	urlMu.Lock()
	urlJobs++
	urlFb = true
	urlMu.Unlock()
}

func postURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

	s := bufio.NewScanner(r.Body)
	for s.Scan() {
		go func(text string) {
			url := strings.TrimSpace(text)
			if !isLegitURL(url) {
				return
			}
			if err := scheduler.PushURL(r.Context(), url); err != nil {
				stderr.Printf("pushURL (%s): %s", url, err)
			}
		}(s.Text())
	}
	if err := s.Err(); err != nil {
		stderr.Println("pushURL:", err)
		w.WriteHeader(400)
		return
	}

	select {
	case urlWait <- struct{}{}:
	default:
	}

	w.WriteHeader(204) // 204 No Content

	htmlMu.Lock()
	htmlJobs--
	if htmlJobs < 0 {
		htmlJobs = 0
	}
	htmlMu.Unlock()
}

func getHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

REFETCH:
	html, isDrained, err := scheduler.PopHTML(r.Context())
	if isDrained {
		if isOwari() {
			stdout.Println("OWATTA! (from getHTML)")
			w.WriteHeader(410)
			return
		}
		stdout.Println("HTML temporarily drained, waiting.")
		select {
		case <-htmlWait:
			goto REFETCH
		case <-r.Context().Done():
			return
		}
	}
	if err != nil {
		stderr.Println("getHTML:", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(html)))
	w.WriteHeader(200)
	io.WriteString(w, html)

	htmlMu.Lock()
	htmlJobs++
	htmlFb = true
	htmlMu.Unlock()
}

func postHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		stderr.Println("postHTML:", err)
		w.WriteHeader(400)
		return
	}

	// 这里做 trim 的成本有点高，直接跳过了
	if buf.Len() == 0 {
		w.WriteHeader(400)
		return
	}

	html := buf.String()
	if err := scheduler.PushHTML(r.Context(), html); err != nil {
		stderr.Println("postHTML:", err)
		w.WriteHeader(400)
		return
	}

	select {
	case htmlWait <- struct{}{}:
	default:
	}

	w.WriteHeader(204)

	urlMu.Lock()
	urlJobs--
	if urlJobs < 0 {
		urlJobs = 0
	}
	urlMu.Unlock()
}
