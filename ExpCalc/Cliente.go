package main

import (
	"bufio"
	"fmt"
	"net"
	//"github.com/stdi0/calc"
	"os"
)

type Expression struct{
	a string
}

func main() {
	// estabeler conexão TCP socket
	con, _ := net.Dial("tcp", "127.0.0.1:8080")
	for {
		// criar variável para ler expressão do cliente
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Informe a expressão: (e.g., 8-4+6)")
		text, _ := reader.ReadString('\n') //lendo a expressão e colocando o conteúdo na variável text
		_, _ = fmt.Fprintf(con, text+"\n") //enviar conteúdo da mensagem via a conexão socket
		result, _ := bufio.NewReader(con).ReadString('\n') //variável para ler resposta do servidor
		fmt.Print("result: " + result) //exibir resposta
		os.Exit(1)
	}

}
