package http

import (
	"io"
	"net/http"
)

func delayWriter(r io.ReadCloser, w http.ResponseWriter) (io.ReadCloser, http.ResponseWriter) {
	ch := make(chan struct{})

	return &requestBody{r, ch}, &responseWriter{w, ch}
}

type requestBody struct {
	io.ReadCloser

	ch chan<- struct{}
}

func (r *requestBody) Read(buf []byte) (int, error) {
	n, err := r.ReadCloser.Read(buf)
	if err != nil {
		close(r.ch)
	}

	return n, err
}

type responseWriter struct {
	http.ResponseWriter

	ch <-chan struct{}
}

func (w *responseWriter) Write(buf []byte) (int, error) {
	<-w.ch
	return w.ResponseWriter.Write(buf)
}

func (w *responseWriter) WriteHeader(status int) {
	<-w.ch
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriter) Flush() {
	<-w.ch
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
