# Vigil - Rastreador de Preços Online 👁
[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## 📋 Sobre

O **Vigil** é um rastreador inteligente de preços desenvolvido em [Go](https://go.dev/), projetado para monitorar ofertas online e alertar automaticamente quando o **menor valor histórico** for atingido.

Com suporte para **Amazon**, **Kabum** e **Mercado Livre**, o Vigil ajuda você a economizar, notificando diretamente no **Telegram** quando uma promoção real é detectada.

## 🚀 Funcionalidades

- 🔍 Monitora produtos de múltiplas lojas (Amazon, Kabum, Mercado Livre)
- 📊 Mantém histórico de preços e realiza análises
- 📉 Detecta automaticamente menores preços históricos
- 📱 Envia alertas de promoções via [Telegram Bot API](https://core.telegram.org/bots/api)
- 📄 Importação de produtos via CSV
- 🛠️ Interface CLI completa e intuitiva
- 📋 Visualização em tabela dos preços históricos

## 🛠️ Tecnologias

- [Go](https://go.dev/) - Linguagem de programação principal
- [Colly](https://github.com/gocolly/colly) - Web scraping
- [Telegram Bot API](https://core.telegram.org/bots/api) - Notificações
- [PostgreSQL](https://www.postgresql.org/) - Armazenamento de dados
- [Cobra](https://github.com/spf13/cobra) - Interface CLI
- [tablewriter](https://github.com/olekukonko/tablewriter) - Visualização de dados em tabela

## 🔧 Pré-requisitos

- [Go 1.24+](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- Conexão com a internet

## ⚙️ Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/ViniMendes2515/vigil.git
cd vigil
```

### 2. Configure o .env

Crie o arquivo .env com as seguintes variáveis:

```bash
TELEGRAM_TOKEN=seu-token-do-telegram
TELEGRAM_CHAT_ID=seu-chat-id
DATABASE_URL=sua-url-do-postgres
```

### 3. Compile o projeto
```bash
go build -o vigil
```
### 4. Execute o Vigil
```bash
./vigil
```

### 🐳 Docker
Você pode rodar o Vigil via Docker:

```bash
docker build -t vigil .
docker run -it --rm --env-file .env vigil
```

📘 Uso da CLI

✅ Adicionar um produto
```bash
vigil add --site amazon --url "https://www.amazon.com.br/dp/123456" --name "Meu Produto" --preco 99.90
```
✅ Listar produtos monitorados
```bash
vigil list
```
✅ Ver detalhes de um produto
```bash
vigil show --url "https://www.amazon.com.br/dp/123456"
```
✅ Verifica se esta em promoção, se estiver envia alerta via Telegram
```bash
vigil check
```

✅ Remover um produto
```bash

vigil del --url "https://www.amazon.com.br/dp/123456"
```

🧪 Testes
Rode os testes unitários com:
```bash
go test ./...
```

📄 Licença <br>
Este projeto está licenciado sob a Licença MIT.

✨ Autor
Desenvolvido por Vinicius Mendes
