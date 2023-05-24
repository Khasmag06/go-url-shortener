package gRPC

import (
	"context"
	"errors"
	"fmt"
	"github.com/Khasmag06/go-url-shortener/config"
	"github.com/Khasmag06/go-url-shortener/internal/app/handlers/gRPC/interceptors"
	"github.com/Khasmag06/go-url-shortener/internal/app/shorten"
	"github.com/Khasmag06/go-url-shortener/internal/app/storage"
	pb "github.com/Khasmag06/go-url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type ShortenerServer struct {
	pb.UnimplementedShortenerServer
	cfg  config.Config
	repo storage.Storage
}

func NewShortenerServer(cfg config.Config, repo storage.Storage) *ShortenerServer {
	return &ShortenerServer{cfg: cfg, repo: repo}
}

func (s *ShortenerServer) Run() {

	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Println(err)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(interceptors.CreateAccessToken))
	pb.RegisterShortenerServer(srv, s)

	fmt.Println("Сервер gRPC начал работу")
	go func() {
		if err := srv.Serve(listen); err != nil {
			log.Println(err)
		}
	}()
}

func (s *ShortenerServer) AddShortURL(
	ctx context.Context, in *pb.AddShortRequest) (*pb.AddShortResponse, error) {
	var response pb.AddShortResponse
	md, _ := metadata.FromIncomingContext(ctx)
	userID := md.Get("token")[0]
	short := shorten.URLShorten()
	shortURL := storage.ShortURL{ID: short, OriginalURL: in.OriginalURL}
	err := s.repo.AddShortURL(userID, &shortURL)
	if err != nil && errors.Is(err, storage.ErrExistsURL) {
		short, err = s.repo.GetExistURL(in.OriginalURL)
		if err != nil {
			log.Println(err)
		}
		response.ShortURL = fmt.Sprintf("%s/%s", s.cfg.BaseURL, short)
		return &response, nil
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &response, nil
}
