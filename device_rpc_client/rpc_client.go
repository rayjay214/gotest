package main

import (
	"fmt"
	"log"

	"context"
	"google.golang.org/grpc"
	"gotest/device_rpc_client/proto"
)

func main() {
	conn, err := grpc.Dial("114.215.190.173:40051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	sendCmd(conn)
	//openShortRecord(conn)
	//vorRecordSwitch(conn)
}

func sendCmd(conn *grpc.ClientConn) {
	client := proto.NewDeviceServiceClient(conn)

	req := &proto.SendCmdRequest{
		Imei: 13320465357,
		//Content: "SL DP114.215.190.173#8881#",
		Content: "SL VERSION",
	}
	fmt.Println(req.String())

	resp, err := client.SendCmd(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling SendCmd: %v", err)
	}

	log.Printf("Response from server: %v", resp)
}

func openShortRecord(conn *grpc.ClientConn) {
	client := proto.NewDeviceServiceClient(conn)

	req := &proto.OpenShortRecordRequest{
		Imei:    26191697155,
		Seconds: 10,
	}
	fmt.Println(req.String())
	resp, err := client.OpenShortRecord(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling SendCmd: %v", err)
	}

	log.Printf("Response from server: %v", resp)
}

func vorRecordSwitch(conn *grpc.ClientConn) {
	client := proto.NewDeviceServiceClient(conn)

	req := &proto.VorRecordSwitchRequest{
		//Imei: 13320465357,
		Imei:   58231222504,
		Switch: 1,
	}
	resp, err := client.VorRecordSwitch(context.Background(), req)

	if err != nil {
		log.Fatalf("error calling SendCmd: %v", err)
	}

	log.Printf("Response from server: %v", resp)
}
