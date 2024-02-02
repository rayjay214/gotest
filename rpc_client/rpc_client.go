package main

import (
    "context"
    "log"

    "google.golang.org/grpc"
    "gotest/rpc_client/proto"
)

func main() {
    conn, err := grpc.Dial("114.215.190.173:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("could not connect: %v", err)
    }
    defer conn.Close()

    client := proto.NewIpcServiceClient(conn)

    req := &proto.SendStunAddrRequest{
        //Uid:  "0102030405060708",
        Uid:  "123456789111111",
        Ip:   "113.110.215.16",
        Port: 22345,
    }
    resp, err := client.SendStunAddr(context.Background(), req)
    if err != nil {
        log.Fatalf("error calling SayHello: %v", err)
    }

    log.Printf("Response from server: %v", resp)
}
