package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Сколько в среднем спим при эмуляции работы
const AvgSleep = 50

type Timing struct {
	Count    int
	Duration time.Duration
}

type ctxTimings struct {
	sync.Mutex
	Data map[string]*Timing
}

// Линтер ругается если используем базовые типы в Value контекста, поэтому заводим кастомный
type key int

const timingsKey key = 1

func logContextTimings(ctx context.Context, path string, start time.Time) {
	// Получаем тайминги из контекста.
	// Поскольку там пустой интерфейс, то нам надо преобразовать к нужному типу:
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}

	realTotalElapsed := time.Since(start)

	buf := bytes.NewBufferString(path)
	var totalElapsedFromTimings time.Duration
	for timing, value := range timings.Data {
		totalElapsedFromTimings += value.Duration
		buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
	}

	buf.WriteString(fmt.Sprintf("\n\treal total: %s", realTotalElapsed))
	buf.WriteString(fmt.Sprintf("\n\ttracked total: %s", totalElapsedFromTimings))
	buf.WriteString(fmt.Sprintf("\n\tunknown: %s", realTotalElapsed-totalElapsedFromTimings))

	fmt.Println(buf.String())
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx != r.Context(), r.Context() has no ctxTimings
		ctx := context.WithValue(
			r.Context(),
			timingsKey,
			&ctxTimings{
				Data: make(map[string]*Timing),
			},
		)

		defer logContextTimings(ctx, r.URL.Path, time.Now())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func trackContextTimings(ctx context.Context, metricName string, start time.Time) {
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}

	elapsedTime := time.Since(start)

	// Лочимся на случай конкурентной записи в мапу:
	timings.Lock()
	defer timings.Unlock()

	// Если метрики ещё нет - мы её создадим, если есть - допишем в существующую:
	if metric, metricExists := timings.Data[metricName]; metricExists {
		metric.Count++
		metric.Duration += elapsedTime
	} else {
		timings.Data[metricName] = &Timing{
			Count:    1,
			Duration: elapsedTime,
		}
	}
}

func emulateWork(ctx context.Context, workName string) {
	defer trackContextTimings(ctx, workName, time.Now())

	rnd := time.Duration(rand.Intn(AvgSleep))
	time.Sleep(time.Millisecond * rnd)
}

func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	emulateWork(ctx, "checkCache")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")

	time.Sleep(10 * time.Millisecond) // Не отслеживаемая в таймингах трата времени

	emulateWork(ctx, "loadSidebar")
	emulateWork(ctx, "loadComments")

	fmt.Fprintln(w, "Request done")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", loadPostsHandle)

	siteHandler := timingMiddleware(mux)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", siteHandler)
}
