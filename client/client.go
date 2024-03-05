package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func criarArquivo(valor float64) error {
	file, err := os.Create("cotacao.txt")

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %.2f", valor))

	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		log.Printf("Não foi possível efetuar a conexão: %v\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Ocorreu um erro ao ao receber a requisição: %v\n", err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Printf("Ocorreu um erro ao ler a requisição: %v\n", err)
		return
	}

	var data cotacao

	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("Ocorreu um erro ao converter os dados para json: %v\n", err)
		return
	}

	valor_bid, err := strconv.ParseFloat(data.Bid, 64)

	if err != nil {
		log.Printf("Ocorreu um erro ao converter uma string para float: %v\n", err)
		return
	}

	err = criarArquivo(valor_bid)

	if err != nil {
		log.Printf("Ocorreu um erro ao criar o arquivo: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Valor atual do Dolar US$: %.2f\n", valor_bid)
}
