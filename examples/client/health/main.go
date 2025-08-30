package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr  = flag.String("addr", "localhost:6666", "The grpc server address to connect to.") //grpc的服务地址
	limit = flag.Int64("limit", 10, "Limit to list users.")                                 //限制列出用户的数量
)

func main() {
	flag.Parse()
	//建立与grpc的连接
	//grcp.NewClient 用于建立客户端与服务端的连接
	//grpc.WithTransportCredentials(insecure.NewCredentials()) 表示使用不安全的传输（即不使用 TLS）
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	//创建MiniBlog客户端
	client := apiv1.NewMiniBlogClient(conn)
	// 设置上下文，带有 3 秒的超时时间
	// context.WithTimeout 用于设置调用的超时时间，防止请求无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() //在函数结束后取消上下文，释放资源
	//调用MiniBlog的Health方法，检查服务状态
	resp, err := client.Healthz(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to call healthz: %v", err)
	}
	// 将返回的响应数据转换为 JSON 格式
	jsondata, _ := json.Marshal(resp)
	fmt.Println(string(jsondata))

}
