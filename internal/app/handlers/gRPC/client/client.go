package main

import (
	"context"
	"github.com/Khasmag06/go-url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := proto.NewShortenerClient(conn)

	md := metadata.New(map[string]string{"token": "123456"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = c.AddShortURL(ctx, &proto.AddShortRequest{OriginalURL: "https://www.google10.com/"})
	if err != nil {
		log.Fatal(err)
	}

}
