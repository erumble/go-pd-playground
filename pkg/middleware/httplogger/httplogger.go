package httplogger

import (
	"log"
	"net/http"
	"time"
)

type interceptResponseWriter struct {
	http.ResponseWriter
	HTTPStatus   int
	ResponseSize int
}

func (w *interceptResponseWriter) WriteHeader(status int) {
	w.HTTPStatus = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *interceptResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *interceptResponseWriter) Write(b []byte) (int, error) {
	if w.HTTPStatus == 0 {
		w.HTTPStatus = http.StatusOK
	}

	w.ResponseSize = len(b)
	return w.ResponseWriter.Write(b)
}

//HTTPLogger logs http requests in Apache Common Log Format
func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		interceptWriter := interceptResponseWriter{w, 0, 0}

		next.ServeHTTP(&interceptWriter, r)
		log.Printf("HTTP - %s - - %s \"%s %s %s\" %d %d %s %dus\n",
			r.RemoteAddr,
			t.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			interceptWriter.HTTPStatus,
			interceptWriter.ResponseSize,
			r.UserAgent(),
			time.Since(t),
		)
	})
}
