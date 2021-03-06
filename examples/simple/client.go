package main

import (
	"fmt"
	"os"

	"nanomsg.org/go-mangos"
	"nanomsg.org/go-mangos/protocol/req"
	"nanomsg.org/go-mangos/transport/ipc"
	"nanomsg.org/go-mangos/transport/tcp"

	"github.com/golang/protobuf/proto"
	pb "github.com/lucidity-dev/bulletin/protobuf"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func main() {
	url := "tcp://127.0.0.1:40899"
	var err error
	var data []byte
	var sock mangos.Socket
	var msg []byte
	
	sock, err = req.NewSocket()
	if err != nil {
		die("error creating socket: %s", err.Error())
	}

	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		die("cant dial on socket: %s", err)
	}

	request := &pb.Message{
		Cmd: pb.Message_REGISTER,
		Args: "test1",
	}

	data, err = proto.Marshal(request)

	if err != nil {
		die("error creating protobuf: %s", err)
	}

	if err = sock.Send(data); err != nil {
		die("error sending message: %s", err)
	}

	if msg, err = sock.Recv(); err != nil {
		die("didn't receive any data: %s", err)
	}
	//fmt.Println(string(msg))
	var body pb.Topic
	proto.Unmarshal(msg, &body)

	if body.Err == "" {
		fmt.Printf("Topic: %s -> URL: %s", string(body.Name), string(body.Url))
	} else {
		fmt.Println(body.Err)
	}

	request = &pb.Message{
		Cmd: pb.Message_GET,
		Args: "test1",
	}

	data, err = proto.Marshal(request)

	if err != nil {
		die("error creating protobuf: %s", err)
	}

	if err = sock.Send(data); err != nil {
		die("error sending message: %s", err)
	}

	if msg, err = sock.Recv(); err != nil {
		die("didn't receive any data: %s", err)
	}

	proto.Unmarshal(msg, &body)

	if body.Err == "" {
		fmt.Printf("Topic: %s -> URL: %s", string(body.Name), string(body.Url))
	} else {
		fmt.Println(body.Err)
	}

	sock.Close()
}
