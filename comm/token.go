package comm

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

//func GetClientInterceptor() grpc.StreamClientInterceptor {
//
//	var clientinterceptor grpc.StreamClientInterceptor
//
//	clientinterceptor = func(
//		ctx context.Context,
//		desc *grpc.StreamDesc,
//		cc *grpc.ClientConn,
//		method string,
//		streamer grpc.Streamer,
//		opts ...grpc.CallOption,
//	) (grpc.ClientStream, error) {
//
//		return streamer(ctx, desc, cc, method, opts...)
//	}
//	return clientinterceptor
//}

func GetServerInterceptor() grpc.StreamServerInterceptor {
	var serverinterceptor grpc.StreamServerInterceptor

	serverinterceptor = func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := Auth(ss.Context())
		if err != nil {
			//return status.Errorf(codes.Unauthenticated, err.Error())
			return err
		}
		return handler(srv, ss)
	}
	return serverinterceptor
}

func Auth(ctx context.Context) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var user string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	//if user != "joy" && user != "jack" && user != "tom" {
	if user == "" {
		return status.Errorf(codes.Unauthenticated, "token 不合法")
	}
	return nil

}
