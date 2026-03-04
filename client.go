package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type RespostaCotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("Erro na requisição: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Server retornou erro (%d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Erro ao ler resposta: %v", err)
	}

	var data RespostaCotacao
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Erro ao fazer Unmarshal no JSON: %v", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Erro ao criar o arquivo: %v", err)
	}

	defer file.Close()

	content := fmt.Sprintf("Dólar: %s", data.Bid)
	_, err = file.WriteString(content)

	if err != nil {
		log.Fatalf("Erro ao salvar no arquivo: %v", err)
	}

	fmt.Println("Cotação salva com sucesso em cotacao.txt")
}
