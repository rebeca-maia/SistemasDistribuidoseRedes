package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	sc "text/scanner"
)

type Token uint

type Expression struct{
	a string
}

const (
	ILLEGAL Token = iota

	NUMBER
	ADDITION
	SUBTRACTION
	MULTIPLICATION
	DIVISION

	LBRACKET
	RBRACKET
)

type Lexeme struct {
	Token Token
	Value string
}

// passando a entrada em string para Tokens
func Lexer(r io.Reader) (out []Lexeme) {
	var s sc.Scanner
	s.Init(r)
	s.Mode = sc.ScanInts

	var tok rune
	for tok != sc.EOF {
		tok = s.Scan()
		mask := 1 << -uint(tok)

		var lexeme Lexeme
		if mask == sc.ScanInts {
			lexeme = Lexeme{NUMBER, s.TokenText()}
			out = append(out, lexeme)
			continue
		} else if tok == sc.EOF {
			break
		}

		switch string(tok) {
		case "+":
			lexeme = Lexeme{ADDITION, s.TokenText()}
		case "-":
			lexeme = Lexeme{SUBTRACTION, s.TokenText()}
		case "*":
			lexeme = Lexeme{MULTIPLICATION, s.TokenText()}
		case "/":
			lexeme = Lexeme{DIVISION, s.TokenText()}
		case "(":
			lexeme = Lexeme{LBRACKET, s.TokenText()}
		case ")":
			lexeme = Lexeme{RBRACKET, s.TokenText()}
		default:
			lexeme = Lexeme{ILLEGAL, s.TokenText()}
		}
		out = append(out, lexeme)
	}

	return
}

// convertendo os Tokens para uma sequência na Notação Polonesa Reversa
func Parser(in []Lexeme) (rpn []Lexeme, err error) {
	stack := make([]Lexeme, 0)
	var top Lexeme
	for num, lex := range in {
		if len(stack) > 0 {
			top = stack[len(stack)-1]
		}

		switch {
		case lex.Token == ILLEGAL:
			return nil, fmt.Errorf("illegal token detected")
		case lex.Token == NUMBER:
			rpn = append(rpn, lex)
		case len(stack) == 0:
			if len(rpn) == 0 && lex.Token == SUBTRACTION {
				rpn = append(rpn, Lexeme{NUMBER, "0"})
			}
			stack = append(stack, lex)
		case lex.Token == LBRACKET:
			stack = append(stack, lex)
		case lex.Token == RBRACKET:
			for i := len(stack) - 1; i >= 0; i-- {
				top := stack[i]
				stack = stack[:i]
				if top.Token == RBRACKET {
					return nil, fmt.Errorf("unpaired bracket")
				} else if top.Token == LBRACKET {
					break
				} else if i == 0 {
					return nil, fmt.Errorf("brackets not matched")
				}
				rpn = append(rpn, top)
			}
		case lex.Token == ADDITION || lex.Token == SUBTRACTION:
			prevLex := in[num-1]
			if top.Token != LBRACKET && top.Token != RBRACKET {
				rpn = append(rpn, top)
				stack = stack[:len(stack)-1]
			} else if lex.Token == SUBTRACTION && prevLex.Token == LBRACKET {
				rpn = append(rpn, Lexeme{NUMBER, "0"})
			}
			stack = append(stack, lex)
		case (lex.Token == MULTIPLICATION || lex.Token == DIVISION) &&
			(top.Token == MULTIPLICATION || top.Token == DIVISION):
			rpn = append(rpn, top)
			stack = stack[:len(stack)-1]
			stack = append(stack, lex)
		default:
			stack = append(stack, lex)
		}
	}

	for i := len(stack) - 1; i >= 0; i-- {
		top := stack[i]
		stack = stack[:i]
		if top.Token != LBRACKET && top.Token != RBRACKET {
			rpn = append(rpn, top)
		} else {
			return nil, fmt.Errorf("brackets not matched")
		}
	}

	return
}



func main() {
	ln, _ := net.Listen("tcp", ":8080") //"escutar" mensagens vindas por essa porta
	con, _ := ln.Accept() //aceitar qualquer conexão que for requisitada pelo cliente
	for {
		message, _ := bufio.NewReader(con).ReadString('\n') //lendo expressão vindo da conexão
		in := bufio.NewReader(strings.NewReader(message)) //transformando a string em Reader

		lexems := Lexer(in) //criando um objeto Lexer a partir do objeto Reader
		if lexems == nil { //exibindo o erro
			fmt.Println("error in Lexer")
		}
		rpn, err := Parser(lexems) //criando um objeito Parser a partir do Lexer
		if err != nil {
			fmt.Println(err)
		}

		stack := make([]int, 0) //criando uma pilha de tamanho 0
		for _, lex := range rpn { //iterando cada valor da expressão
			switch {
			case lex.Token == NUMBER:
				intTok, _ := strconv.Atoi(lex.Value) //convertendo valores string para int
				stack = append(stack, intTok) //adicionado valores na pilha
			case len(stack) < 2:
				fmt.Errorf("invalid expression")
			case lex.Token == ADDITION: //realizando a operação de adição com os tokens e colocando o resultado na pilha
				res := stack[len(stack)-2] + stack[len(stack)-1]
				stack = stack[:len(stack)-2]
				stack = append(stack, res)
			case lex.Token == SUBTRACTION: //realizando a operação de subtração com os tokens e colocando o resultado na pilha
				res := stack[len(stack)-2] - stack[len(stack)-1]
				stack = stack[:len(stack)-2]
				stack = append(stack, res)
			case lex.Token == MULTIPLICATION: //realizando a operação de multiplicação com os tokens e colocando o resultado na pilha
				res := stack[len(stack)-2] * stack[len(stack)-1]
				stack = stack[:len(stack)-2]
				stack = append(stack, res)
			case lex.Token == DIVISION: //realizando a operação de divisão com os tokens e colocando o resultado na pilha
				res := stack[len(stack)-2] / stack[len(stack)-1]
				stack = stack[:len(stack)-2]
				stack = append(stack, res)
			}
		}

		res := strconv.Itoa(stack[0]) //convertendo o resultado para string
		//fmt.Print("Resultado:")
		//fmt.Fprint(output,res)
		con.Write([]byte(res)) //escrevendo a resposta na conexão
		os.Exit(1)
	}}


