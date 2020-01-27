package util

import "net/http"

type HttpWriter struct {
	http.ResponseWriter
	Status      int
	wroteHeader bool
}

func (w *HttpWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.Status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *HttpWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
