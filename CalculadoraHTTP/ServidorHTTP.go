package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Operation struct {
	Operand1 string
	Operator string
	Operand2 string
}

func calculate(operand1, operator, operand2 string) string {

	n1, _ := strconv.Atoi(operand1)
	n2, _ := strconv.Atoi(operand2)

	fmt.Println(n1 + n2)
	var res int
	if operator == "+" {
		res = n1 + n2
	} else if operator == "-" {
		res = n1 - n2
	} else if operator == "*" {
		res = n1 * n2
	} else if operator == "/" {
		if(n2==0){ fmt.Println(" Erro! Divisão por zero")}
		res = n1 / n2
	} else {
		fmt.Println("Error in input")
		return ""
	}

	return strconv.Itoa(res)

}

func prepareRespond(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var t Operation
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t.Operand1 + t.Operator + t.Operand2)
	_, _ = w.Write([]byte(calculate(t.Operand1, t.Operator, t.Operand2)))

}
func main() {
	http.HandleFunc("/", prepareRespond)
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println("Error Occurred! ",err)
	}
}