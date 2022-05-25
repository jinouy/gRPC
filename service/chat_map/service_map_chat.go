package main

import (
	"context"
	"fmt"
	"gRPC_User/client/auth"
	pb "gRPC_User/proto/chat"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"sync"
	"time"
)

type ConnectPool struct {
	Map  map[string]pb.OnLineChat_SayHiServer
	Lock *sync.RWMutex
}

func NewConcurMap() *ConnectPool {
	return &ConnectPool{
		Map:  make(map[string]pb.OnLineChat_SayHiServer, 10),
		Lock: &sync.RWMutex{},
	}
}

func (p *ConnectPool) Get(name string) pb.OnLineChat_SayHiServer {
	p.Lock.RLock()

	defer p.Lock.RUnlock()

	if stream, ok := p.Map[name]; ok {
		return stream
	} else {
		return nil
	}
}

func (p *ConnectPool) Add(name string, stream pb.OnLineChat_SayHiServer) {
	p.Lock.Lock()

	defer p.Lock.Unlock()

	p.Map[name] = stream.(pb.OnLineChat_SayHiServer)
}

func (p *ConnectPool) Del(name string) {
	p.Lock.Lock()

	defer p.Lock.Unlock()

	delete(p.Map, name)
}

func (p *ConnectPool) BroadCast(from, message string) bool {

	defer p.Lock.RUnlock()
	p.Lock.RLock()
	log.Printf("message: %s\n", message)
	for username, stream_i := range p.Map {
		stream := stream_i
		if username != from {
			stream.Send(&pb.HiReply{
				Message: message,
				TS:      &timestamp.Timestamp{Seconds: time.Now().Unix()},
			})
		}
	}
	return true
}

type Service struct {
	pb.UnimplementedOnLineChatServer
}

var connect_pool *ConnectPool

func (s *Service) SayHi(stream pb.OnLineChat_SayHiServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unimplemented, "no token")
	}

	var User string
	if val, ok := md["user"]; ok {
		User = val[0]
	} else {
		return status.Errorf(codes.Unimplemented, "no token")
	}
	username := User

	if connect_pool.Get(username) != nil {
		return status.Errorf(codes.Unimplemented, "名字已经存在")
	} else { // 连接成功
		connect_pool.Add(username, stream)
		stream.Send(&pb.HiReply{
			Message: fmt.Sprintf("连接成功!"),
		})
	}

	go func() {
		// 阻塞住，等待断开连接的时候触发
		<-stream.Context().Done()
		connect_pool.Del(username)
		connect_pool.BroadCast(username, fmt.Sprintf("%s 离开房间", username))
	}()

	// 广播，xxxx进入了聊天室直播间
	connect_pool.BroadCast(username, fmt.Sprintf("欢迎 %s!", username))

	//  阻塞接收 该用户后续传来的消息
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		connect_pool.BroadCast(username, fmt.Sprintf("%s: %s", username, req.Message))
	}

}

func main() {

	var opts []grpc.ServerOption

	// 注册一个拦截器
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	connect_pool = NewConcurMap()
	// 监听一个 地址:端口
	address, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Printf("Failed to listen: %v", err)
		return
	}

	// 实例化grpc Server，并开启拦截器
	ser := grpc.NewServer(opts...)
	pb.RegisterOnLineChatServer(ser, &Service{}) //必须实现protoes中定义的方法，不然这里无法通过检测

	// 启动服务
	if err := ser.Serve(address); err != nil {
		log.Printf("Failed to start: %v", err)
		return
	}
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// 进行认证
	name := auth.InputName()
	err := myAuth(ctx, name)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)

}

// 认证token
func myAuth(ctx context.Context, name string) error {

	// md 是一个map[string][]string类型
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unimplemented, "no token")
	}
	var User string

	if val, ok := md["user"]; ok {
		User = val[0]
	}

	if User != name {
		return status.Errorf(codes.Unimplemented, "token invalide: user=%s", User)
	}
	return nil
}
