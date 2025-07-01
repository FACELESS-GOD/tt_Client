package main

import (
	"context"
	"io"
	"log"
	"time"

	tt_git "github.com/FACELESS-GOD/tt_Client.git/Package/ProtocGenerated/ProtocBuff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	clientConn := tt_git.NewTestGRPCServiceClient(conn)

	customCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	inst, err := clientConn.GetData(customCtx, &tt_git.DataReceipt{})

	if err != nil {
		panic(err)
	}

	dtchan := make(chan bool)
	var items []*tt_git.Item

	go func() {
		for {
			resp, err := inst.Recv()

			if err == io.EOF {
				dtchan <- true
				return
			} else if err != nil {
				panic(err)
			} else {
				items = resp.Items
				log.Print("REcieved")
			}

		}
	}()

	<-dtchan

	r, err := clientConn.AddData(customCtx, &tt_git.SendData{Items: items})

	log.Println(r)

}
