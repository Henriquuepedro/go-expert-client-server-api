package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type ConexaoBanco struct {
	DB *sql.DB
}

type Cotacao struct {
	Bid        string `json:"bid"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ConsultaCotacao struct {
	USDBRL Cotacao `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema := "CREATE TABLE IF NOT EXISTS cotacoes (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"bid TEXT, " +
		"code TEXT, " +
		"codein TEXT, " +
		"name TEXT, " +
		"high TEXT, " +
		"low TEXT, " +
		"varBid TEXT, " +
		"pctChange TEXT, " +
		"ask TEXT, " +
		"timestamp INTEGER, " +
		"create_date DATETIME" +
		");"

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	conexao := &ConexaoBanco{DB: db}

	http.HandleFunc("/cotacao", conexao.getCotacaoHandler)

	log.Println("Iniciado na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func (s *ConexaoBanco) getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	cotacao, err := getCotacao()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = saveToDB(s.DB, cotacao)
	if err != nil {
		log.Printf("Erro ao salvar no banco: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}

func getCotacao() (*Cotacao, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data ConsultaCotacao

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return &data.USDBRL, nil
}

func saveToDB(db *sql.DB, cotacao *Cotacao) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, `INSERT INTO cotacoes (
		bid, code, codein, name, high, low, varBid, pctChange, ask, timestamp, create_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cotacao.Bid,
		cotacao.Code,
		cotacao.Codein,
		cotacao.Name,
		cotacao.High,
		cotacao.Low,
		cotacao.VarBid,
		cotacao.PctChange,
		cotacao.Ask,
		cotacao.Timestamp,
		cotacao.CreateDate,
	)

	if err != nil {
		return err
	}

	return nil
}
