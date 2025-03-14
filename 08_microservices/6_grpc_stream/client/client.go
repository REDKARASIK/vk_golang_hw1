package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc"

	"gitlab.vk-golang.ru/vk-golang/lectures/08_microservices/6_grpc_stream/translit"
)

func main() {
	grpcConn, err := grpc.Dial(
		"127.0.0.1:8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grpcConn.Close()

	tr := translit.NewTransliterationClient(grpcConn)

	ctx := context.Background()
	stream, _ := tr.EnRu(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			outWord, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("\tstream closed")
				return
			} else if err != nil {
				fmt.Println("\terror happened", err)
				return
			}
			fmt.Println(" <-", outWord.Word)
		}
	}(wg)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_ = stream.Send(&translit.Word{
			Word: scanner.Text(),
		})
	}
	stream.CloseSend()

	wg.Wait()

}
