package main

import (
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

const (
	idLength        = 8
	requestIDHeader = "X-Request-ID"
	charset         = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type AccessLogger struct {
	StdLogger    *log.Logger
	StdSlogger   *slog.Logger
	ZapLogger    *zap.SugaredLogger
	LogrusLogger *logrus.Entry
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := generateID()
		w.Header().Add(requestIDHeader, reqID)

		next.ServeHTTP(w, r)

		fmt.Printf(
			"FMT [%s] %s %s: %s %s\n",
			reqID, r.Method, r.RemoteAddr, r.URL.Path, time.Since(start),
		)

		log.Printf(
			"LOG [%s] %s %s: %s %s\n",
			reqID, r.Method, r.RemoteAddr, r.URL.Path, time.Since(start),
		)

		ac.StdLogger.Printf(
			"[%s] %s %s: %s %s\n",
			reqID, r.Method, r.RemoteAddr, r.URL.Path, time.Since(start),
		)

		ac.StdSlogger.Info(
			"new api request",
			"request_id", reqID,
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL,
			"work_time", time.Since(start),
			"token", "MY_VERY_SECRET_TOKEN", // Светим в логах секрет
		)

		ac.ZapLogger.Info(
			r.URL.Path,
			zap.String("request_id", reqID),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"request_id":  reqID,
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

func generateID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "\n\tHello world!\n\n")
}

func main() {
	addr := "localhost"
	port := 8080

	// std fmt
	fmt.Printf("STD fmt starting server at %s:%d\n", addr, port)

	// std log
	log.Printf("STD log starting server at %s:%d\n", addr, port)

	// std slog - по умолчанию использует стандартный логгер "log"
	slog.Info("STD slog starting server", "addr", addr, "port", port) // Msg and key-value pairs

	// std slog с кастомным логгером
	removeSecrets := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == "token" {
			a.Value = slog.StringValue("*****")
		}

		return a
	}
	// sLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	sLogger := slog.New(slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: removeSecrets},
	))
	sLogger = sLogger.With("logger", "slog")
	sLogger.Info("custom slog starting server", "addr", addr, "port", port)

	// zap
	// у zap-а нет логгера по-умолчанию
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	zapLogger.Info(
		"starting server",
		zap.String("logger", "ZAP"),
		zap.String("host", addr),
		zap.Int("port", port),
	)

	// logrus
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   addr,
		"port":   port,
	}).Info("Starting server")

	accessLogOut := new(AccessLogger)

	// std
	accessLogOut.StdLogger = log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)

	// std slog
	accessLogOut.StdSlogger = sLogger

	// zap
	zapWithSugar := zapLogger.Sugar().With(
		zap.String("mode", "[access_log]"),
		zap.String("logger", "ZAP"),
	)
	accessLogOut.ZapLogger = zapWithSugar

	// logrus
	logrusLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	accessLogOut.LogrusLogger = logrusLogger

	// server stuff
	siteMux := http.NewServeMux()
	siteMux.HandleFunc("/", mainPage)
	siteHandler := accessLogOut.accessLogMiddleware(siteMux)
	http.ListenAndServe(":8080", siteHandler)
}
