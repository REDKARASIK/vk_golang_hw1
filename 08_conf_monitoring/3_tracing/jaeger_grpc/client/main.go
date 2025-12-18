package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"google.golang.org/grpc"

	"gitlab.vk-golang.ru/vk-golang/lectures/08_conf_monitoring/3_tracing/jaeger_grpc/session"
)

var loginFormTmpl = []byte(`
<html>
	<body>
	<form action="/login" method="post">
		Login: <input type="text" name="login">
		Password: <input type="password" name="password">
		<input type="submit" value="Login">
	</form>
	</body>
</html>
`)

var (
	sessManager session.AuthCheckerClient
)

func checkSession(r *http.Request) (*session.Session, error) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	sess, err := sessManager.Check(
		context.Background(),
		&session.SessionID{
			ID: cookieSessionID.Value,
		})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func innerPage(w http.ResponseWriter, r *http.Request) {
	sess, err := checkSession(r)
	if err != nil || sess == nil {
		w.Write(loginFormTmpl)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "Welcome, "+sess.Login+" <br />")
	fmt.Fprintln(w, "Session ua: "+sess.Useragent+" <br />")
	fmt.Fprintln(w, `<a href="/logout">logout</a>`)
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	inputLogin := r.FormValue("login")
	expiration := time.Now().Add(365 * 24 * time.Hour)

	sess, err := sessManager.Create(
		context.Background(),
		&session.Session{
			Login:     inputLogin,
			Useragent: r.UserAgent(),
		})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   sess.ID,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	ctx := context.Background()
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return
	}

	// Создаем resource с информацией о сервисе
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("client"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Создаем провайдер трассировки
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Устанавливаем глобальный провайдер трассировки
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	// Создаем gRPC соединение с OpenTelemetry interceptor'ом
	grcpConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	sessManager = session.NewAuthCheckerClient(grcpConn)

	http.HandleFunc("/", innerPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	cookieSessionID, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sessManager.Delete(
		context.Background(),
		&session.SessionID{
			ID: cookieSessionID.Value,
		})

	cookieSessionID.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookieSessionID)

	http.Redirect(w, r, "/", http.StatusFound)
}
