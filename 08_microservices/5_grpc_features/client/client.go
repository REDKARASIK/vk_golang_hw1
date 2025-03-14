package main

import (
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"gitlab.vk-golang.ru/vk-golang/lectures/08_microservices/4_grpc/session"
)

// {Interceprtor} | serialization | -> request

func timingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Printf(`--
	call=%v
	req=%#v
	reply=%#v
	time=%v
	err=%v
`, method, req, reply, time.Since(start), err)
	return err
}

type tokenAuth struct {
	Token string
}

func (t *tokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"access-token": t.Token,
	}, nil
}

func (c *tokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {

	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithUnaryInterceptor(timingInterceptor),
		grpc.WithPerRPCCredentials(&tokenAuth{"100500"}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConn.Close()

	sessManager := session.NewAuthCheckerClient(grpcConn)

	ctx := context.Background()
	md := metadata.Pairs(
		"api-req-id", "123",
		"subsystem", "cli",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)

	var header, trailer metadata.MD

	// создаем сессию
	sessId, err := sessManager.Create(ctx,
		&session.Session{
			Login:     "rvasily",
			Useragent: "chrome",
		},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	fmt.Println("sessId", sessId, err)
	fmt.Println("header", header)
	fmt.Println("trailer", trailer)

	// проверяем сессию
	sess, err := sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)

	// удаляем сессию
	_, err = sessManager.Delete(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})

	// проверяем еще раз
	sess, err = sessManager.Check(ctx,
		&session.SessionID{
			ID: sessId.ID,
		})
	fmt.Println("sess", sess, err)
}
