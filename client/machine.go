package main

import (
	"context"
	"flag"
	"io"
	"log"
	"microservice/machine"
	"time"

	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:9111", "The server address in the format of host:port")
)

func runExecute(client machine.MachineClient, instructions *machine.InstructionSet) {
	log.Printf("Executing %v", instructions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.Execute(ctx)
	if err != nil {
		log.Fatalf("%v.Execute(ctx) = %v, %v: ", client, stream, err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF")
				close(waitc)
				return
			}

			if err != nil {
				log.Printf("Err: %v", err)
			}

			log.Printf("output: %v", result.GetOutput())
		}
	}()

	for _, instruction := range instructions.GetInstructions() {
		if err := stream.Send(instruction); err != nil {
			log.Fatalf("%v.Send(%v) = %v: ", stream, instruction, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("%v.CloseSend() got error %v, want %v", stream, err, nil)
	}
	<-waitc
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := machine.NewMachineClient(conn)

	// try Execute()
	instructions := []*machine.Instruction{}
	instructions = append(instructions, &machine.Instruction{Operand: 5, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operand: 2, Operator: "PUSH"})
	instructions = append(instructions, &machine.Instruction{Operator: "MUL"})
	instructions = append(instructions, &machine.Instruction{Operator: "FIB"})
	runExecute(client, &machine.InstructionSet{Instructions: instructions})
}
