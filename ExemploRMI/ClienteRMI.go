package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ", os.Args[0], "server:port")
		os.Exit(1)
	}
	service := os.Args[1]+":"+os.Args[2]

	client, err := rpc.Dial("tcp", service)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := Args{17, 8}
	var reply int

	err = client.Call("Arith.Adicao", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Soma: %d+%d=%d\n", args.A, args.B, reply)

	err = client.Call("Arith.Subtracao", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Subtracao: %d-%d=%d\n", args.A, args.B, reply)

	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Multiplicacao: %d*%d=%d\n", args.A, args.B, reply)

	var quot Quotient
	err = client.Call("Arith.Divide", args, &quot)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Divisao: %d/%d=%d resto %d\n", args.A, args.B, quot.Quo, quot.Rem)

}