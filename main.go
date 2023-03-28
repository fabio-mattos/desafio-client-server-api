package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type DollarDtoInput struct {
	USDBRL DollarDto `json:"USDBRL"`
}

type DollarDto struct {
	ID         int    `json:"id"`
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

type DollarDtoOutput struct {
	Bid        string `json:"bid"`
	CreateDate string `json:"create_date"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {

	log.Print(BuscaCotacao())

}

func BuscaCotacao() (*DollarDtoInput, error) {
	resp, error := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var c DollarDtoInput
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	SalvaTxt(c.USDBRL.Bid)
	log.Print("Fabio o Valor do hoje é Dolar: ", c.USDBRL.Bid)

	return &c, nil
}

func SalvaTxt(bid string) {

	//if !(len(bid) > 0) || !(len(dataCotacao) > 0) {

	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Println(err)
	}

	// Salvando a cotação atual no  arquivo "cotacao.txt" no formato: Dólar: {valor}
	tamanho, err := fmt.Fprintln(f, "Dólar: ", bid)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Arquivo criado com sucesso! Tamanho: %d bytes\n", tamanho)
	f.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Arquivo criado com sucesso! ")
	f.Close()
	/*
		} else {
			fmt.Printf("Não foi possível salvar a cotação atual! O bid ou a data da cotação está vazia!")
		}
	*/

}
