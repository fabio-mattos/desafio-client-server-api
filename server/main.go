package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacaodb struct {
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

type Dolar struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", cotaHandler)
	http.ListenAndServe(":8080", nil)

}
func cotaHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	var cotacaoAtual map[string]Cotacaodb
	cotacaoAtual, err := BuscaCotacao(ctx)
	if err != nil {
		panic(fmt.Sprintf("Não foi possíver buscar a cotação %v", err))
	}

	ctx = nil
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Millisecond)

	err = GravandoCotacaoNoBanco(ctx, cotacaoAtual)
	if err != nil {
		panic(fmt.Sprintf("Não foi possível gravar a cotação %v no banco de dados!", err))
	}
	dol := Dolar{Bid: cotacaoAtual["USDBRL"].Bid}

	json.NewEncoder(w).Encode(dol)
}

func BuscaCotacao(c context.Context) (map[string]Cotacaodb, error) {
	var data map[string]Cotacaodb
	req, err := http.NewRequestWithContext(c, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return data, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return data, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
func GravandoCotacaoNoBanco(c context.Context, cota map[string]Cotacaodb) error {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return err
	}

	if c.Err() == context.DeadlineExceeded {
		fmt.Println("Erro: O tempo para conectar com a base de dados foi excedido")
		return err
	}

	defer db.Close()
	cot := cota["USDBRL"]
	sts := `
	CREATE TABLE IF NOT EXISTS cotacao(id INTEGER PRIMARY KEY,code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, varbid TEXT, pctchange TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT);
			INSERT INTO
				cotacao(
					code,
					codein,
					name,
					high,
					low,
					varbid,
					pctchange,
					bid,
					ask,
					timestamp,
					create_date
				)
				VALUES (
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?,
					?
				);
	`

	result, err := db.Exec(sts, cot.Code, cot.Codein, cot.Name, cot.High, cot.Low, cot.VarBid, cot.PctChange, cot.Bid, cot.Ask, cot.Timestamp, cot.CreateDate)

	if err != nil {
		return err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(lastID)
	fmt.Println("Tabela atualizada com sucesso!")
	return nil
}
