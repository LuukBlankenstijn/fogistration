package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry := &customLogEntry{
			request:  r,
			useColor: true,
		}

		// Get request ID if present
		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			entry.requestID = reqID
		}

		// Wrap the response writer to capture status code and bytes written
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()

		// Process the request
		next.ServeHTTP(ww, r)

		// Log after request is complete
		entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
	})
}

type customLogEntry struct {
	request   *http.Request
	requestID string
	useColor  bool
}

func (l *customLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra any) {
	// Build parts of the log message
	var parts []string

	parts = append(parts, fmt.Sprintf("[%sAPI%s]", logging.ApiColor, logging.Reset))

	// Add request ID if present
	if l.requestID != "" {
		parts = append(parts, fmt.Sprintf("%s[%s]%s", logging.Yellow, l.requestID, logging.Reset))
	}

	// Add request method and URL
	scheme := "http"
	if l.request.TLS != nil {
		scheme = "https"
	}

	parts = append(parts, fmt.Sprintf("%s\"%s%s%s%s %s %s\"%s",
		logging.Cyan,
		logging.Magenta, l.request.Method, logging.Reset, logging.Cyan,
		fmt.Sprintf("%s://%s%s", scheme, l.request.Host, l.request.RequestURI),
		l.request.Proto, logging.Reset))

	// Add client address
	parts = append(parts, fmt.Sprintf("from %s", l.request.RemoteAddr))

	// Add status code with color
	var statusColor string
	switch {
	case status < 200:
		statusColor = logging.Blue
	case status < 300:
		statusColor = logging.Green
	case status < 400:
		statusColor = logging.Cyan
	case status < 500:
		statusColor = logging.Yellow
	default:
		statusColor = logging.Red
	}
	parts = append(parts, fmt.Sprintf("- %s%d%s", statusColor, status, logging.Reset))

	// Add bytes written
	parts = append(parts, fmt.Sprintf("%s%dB%s", logging.Cyan, bytes, logging.Reset))

	// Add elapsed time with color
	var elapsedColor string
	if elapsed < 500*time.Millisecond {
		elapsedColor = logging.Green
	} else if elapsed < 5*time.Second {
		elapsedColor = logging.Yellow
	} else {
		elapsedColor = logging.Red
	}
	parts = append(parts, fmt.Sprintf("in %s%s%s", elapsedColor, elapsed, logging.Reset))

	// Determine log level based on status
	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	} else if status >= 400 {
		level = slog.LevelWarn
	}

	// Log the complete message with slog (which will add the timestamp)
	slog.Log(context.Background(), level, strings.Join(parts, " "))
}

// We won't ever call this, but we need it to satisfy the interface
func (l *customLogEntry) Panic(v any, stack []byte) {
	slog.Error("Panic", "error", v, "stack", string(stack))
}
