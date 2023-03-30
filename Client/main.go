package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(fmt.Sprintf("Não foi possivel fazer requisição: %v", err))
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(fmt.Sprintf("Não foi possivel salvar a cotaçao: %v", err))
	}

	var dolar Dolar
	defer f.Close()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("Não foi possivel fazer a requisição %v", err))
	}

	defer res.Body.Close()

	cota, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(fmt.Sprintf("Erro ao pegar dados da requisição: %v", err))
	}

	err = json.Unmarshal(cota, &dolar)
	if err != nil {
		panic(fmt.Sprintf("Não foi possivel gerar o json: %v", err))
	}

	var d string = "Dólar: "
	_, err = f.Write([]byte(d))

	if err != nil {
		panic(fmt.Sprintf("Erro ao abrir o arquivo cotacao.txt: %v", err))
	}
	_, err = f.Write([]byte(dolar.Bid))
	if err != nil {
		panic(fmt.Sprintf("Erro ao escrever no arquivo cotacao.txt: %v", err))
	}
}
