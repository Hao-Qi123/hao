package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"seckil/server/basic/bootstrap"
	_ "seckil/server/basic/bootstrap"
	"seckil/server/handler/service"
	"syscall"
	"time"

	"google.golang.org/grpc"

	__ "seckil/proto"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

func main() {

	if err := bootstrap.ConsulInit(); err != nil {
		log.Fatalf("Consul初始化失败: %v", err)
	}
	log.Println("Consul初始化成功")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	services, err := bootstrap.GetServiceWithLoadBalancer("product")
	if err != nil {
		log.Printf("获取product服务失败: %v", err)
	} else {
		log.Printf("获取到product服务: %s, 地址: %s:%d", services.Service, services.Address, services.Port)
	}

	s := grpc.NewServer()
	__.RegisterProductServer(s, &service.Server{})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")
	if err := bootstrap.ConsulShutdown(); err != nil {
		log.Printf("Consul注销失败: %v", err)
	}
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.GracefulStop()
	log.Println("服务已关闭")
}
