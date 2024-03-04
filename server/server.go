package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chasinfo/go-client-server/server/database"
)

type cotacao struct {
	Usdbrl struct {
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
	} `json:"USDBRL"`
}

type mensagem struct {
	Code     int
	Mensagem string
}

func main() {
	http.HandleFunc("/", getCotacao)
	http.ListenAndServe(":8080", nil)
}

func criarArquivo(data cotacao) error {
	file, err := os.Create("cotacao.txt")

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("VarBid: %s", data.Usdbrl))

	if err != nil {
		return err
	}

	return nil
}

func salvarDadosCotacao(c cotacao) error {

	error := database.AutoMigrate()

	if error != nil {
		return error
	}

	database.DbConnect().Create(&database.TCotacao{
		Code:       c.Usdbrl.Code,
		Codein:     c.Usdbrl.Codein,
		Name:       c.Usdbrl.Name,
		High:       c.Usdbrl.High,
		Low:        c.Usdbrl.Low,
		VarBid:     c.Usdbrl.VarBid,
		PctChange:  c.Usdbrl.PctChange,
		Bid:        c.Usdbrl.Bid,
		Ask:        c.Usdbrl.Ask,
		Timestamp:  c.Usdbrl.Timestamp,
		CreateDate: c.Usdbrl.CreateDate,
	})

	return nil
}

func mensagemHttpResponse(w http.ResponseWriter, rMensagem string, status int) {
	msg := mensagem{
		Code:     http.StatusNotFound,
		Mensagem: rMensagem,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}

func getCotacao(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cotacao" {
		mensagemHttpResponse(w, "Endereço não foi encontrado.", http.StatusNotFound)
		fmt.Fprintf(os.Stderr, "Endereço não foi encontrado\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		mensagemHttpResponse(w, "Não foi possível efetuar a conexão.", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Não foi possível efetuar a conexão: %v\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		mensagemHttpResponse(w, "Ocorreu um erro ao ao receber a requisição.", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Ocorreu um erro ao ao receber a requisição: %v\n", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		mensagemHttpResponse(w, "Ocorreu um erro ao ler a requisição.", http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Ocorreu um erro ao ler a requisição: %v\n", err)
	}

	var data cotacao

	err = json.Unmarshal(body, &data)

	if err != nil {
		mensagemHttpResponse(w, "Ocorreu um erro ao converter os dados em json.", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Ocorreu um erro ao converter os dados em json: %v\n", err)
	}

	err = salvarDadosCotacao(data)

	if err != nil {
		mensagemHttpResponse(w, "Ocorreu um erro ao salvar os dados.", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Ocorreu um erro ao salvar os dados: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data.Usdbrl)
}
