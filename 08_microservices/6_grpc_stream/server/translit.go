package main

import (
	"fmt"
	tr "github.com/essentialkaos/translit/v2"
	"gitlab.vk-golang.ru/vk-golang/lectures/08_microservices/6_grpc_stream/translit"
	"io"
)

type TrServer struct {
	translit.UnimplementedTransliterationServer
	SetSendCallback func(func(string))
}

func (srv *TrServer) EnRu(inStream translit.Transliteration_EnRuServer) error {
	// srv.SetSendCallback(func(s string) {
	// 	out := &translit.Word{
	// 		Word: s,
	// 	}
	// 	inStream.Send(out)
	// })
	// return nil
	// go func() {
	// 	for {
	// 		inStream.Send(&translit.Word{
	// 			Word: "stat",
	// 		})
	// 		time.Sleep(time.Second)
	// 	}
	// }()
	for {
		// time.Sleep(5 * time.Second)
		inWord, err := inStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		out := &translit.Word{
			Word: tr.ISO9A(inWord.Word),
		}
		fmt.Println(inWord.Word, "->", out.Word)
		inStream.Send(out)
	}
	return nil
}

func NewTr() *TrServer {
	return &TrServer{}
}
