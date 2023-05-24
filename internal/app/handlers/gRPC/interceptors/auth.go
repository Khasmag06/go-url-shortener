package interceptors

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func CreateAccessToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) > 0 {
			token = values[0]
		}
	}
	if len(token) == 0 {
		userID := uuid.NewString()
		md.Set("token", userID)
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
	return handler(ctx, req)
}
