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
	htmlWait = make(chan struct{})

	owari bool
	mu    sync.RWMutex
)

func getURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

	mu.RLock()
	o := owari
	mu.RUnlock()
	if o {
		stdout.Println("URL drained! (yet again)")
		w.WriteHeader(410) // 410 Gone
		return
	}

	url, isDrained, err := scheduler.PopURL(r.Context())
	if isDrained {
		stdout.Println("URL drained!")
		mu.Lock()
		owari = true
		mu.Unlock()
		w.WriteHeader(410)
		return
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
}

func postURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

	s := bufio.NewScanner(r.Body)
	for s.Scan() {
		url := strings.TrimSpace(s.Text())
		if url == "" {
			continue
		}
		if err := scheduler.PushURL(r.Context(), url); err != nil {
			stderr.Printf("pushURL (%s): %s", url, err)
			// return
		}
	}
	if err := s.Err(); err != nil {
		stderr.Println("pushURL:", err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(204) // 204 No Content
}

func getHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")

	mu.RLock()
	o := owari
	mu.RUnlock()
	if o {
		stdout.Println("URL drained! (to getHTML)")
		w.WriteHeader(410) // 410 Gone
		return
	}

REFETCH:
	html, isDrained, err := scheduler.PopHTML(r.Context())
	if isDrained {
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
}
