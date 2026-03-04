# Desafio Go Expert: Client-Server API

Este projeto tem a finalidade de aplicar o que foi estudado nas aulas de **Introdução**, **Configuração de ambiente**, **Fundação**, **Pacotes importantes**, **Context** e **Banco de Dados** no curso Go Expert.

---

## 🛠️ Funcionamento do Sistema

### Server.go
O servidor opera na porta `8080` e expõe o endpoint `/cotacao`.
1. **Consumo de API**: Ao receber uma requisição, o servidor consome a API de Câmbio para buscar a cotação USD-BRL.
    - **Timeout**: Máximo de **200ms**.
2. **Persistência**: O servidor registra cada cotação recebida em um banco de dados **SQLite**.
    - **Timeout**: Máximo de **10ms**.
3. **Resposta**: Retorna o resultado da cotação em formato JSON para o cliente.

### Client.go
O cliente realiza uma requisição ao servidor local.
1. **Requisição**: Solicita a cotação ao endpoint `/cotacao`.
    - **Timeout**: Máximo de **300ms**.
2. **Processamento**: Recebe apenas o valor atual do câmbio (`bid`).
3. **Arquivo**: Salva a cotação em um arquivo chamado `cotacao.txt` no formato: `Dólar: {valor}`.

---

## 🚀 Como Executar

### 1. Preparar o Ambiente
Certifique-se de ter o Go instalado. Este projeto utiliza o driver `modernc.org/sqlite` (Pure Go), o que dispensa a necessidade de um compilador C (GCC) no Windows.

No diretório do projeto, execute:
```bash
go mod tidy
```

### 2. Rodar o Servidor
Abra um terminal e execute:
```bash
go run server.go
```

O servidor criará automaticamente o arquivo `cotacoes.db` e a tabela necessária.

### 3. Rodar o Cliente
Em um novo terminal, execute:
```bash
go run client.go
```

O arquivo `cotacao.txt` será gerado na raiz do projeto.

## 📝 Observações sobre os Timeouts
O desafio impõe um limite de 10ms para a persistência no banco de dados. Em alguns sistemas (especialmente Windows com HDs mecânicos), esse tempo pode ser insuficiente, resultando no erro context deadline exceeded logado no console do servidor. Isso demonstra que o gerenciamento de contexto via WithTimeout está funcionando corretamente conforme os requisitos.

## 📂 Estrutura de Arquivos
- `server.go`: Lógica do servidor, API externa e banco de dados.
- `client.go`: Lógica do cliente, requisição e escrita de arquivo.
- `cotacoes.db`: Banco SQLite (gerado automaticamente).
- `cotacao.txt`: Resultado da última cotação (gerado pelo cliente).