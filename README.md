# Weather by CEP API

Sistema desenvolvido em Go que recebe um CEP brasileiro, identifica a cidade e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin).

## Funcionalidades

- Recebe um CEP válido de 8 dígitos
- Busca a localização através da API ViaCEP
- Consulta o clima atual através da API WeatherAPI
- Retorna as temperaturas em três escalas: Celsius, Fahrenheit e Kelvin

## Requisitos

- Go 1.21+
- Docker e Docker Compose (opcional)
- Chave da API WeatherAPI (obtenha em: https://www.weatherapi.com/)

## Configuração

1. Clone o repositório
2. Copie o arquivo de ambiente:
   ```bash
   cp env.example .env
   ```
3. Configure sua chave da API WeatherAPI no arquivo `.env`

## Execução Local

### Usando Go diretamente

```bash
# Configurar variável de ambiente
export WEATHER_API_KEY=sua_chave_aqui

# Executar
go run main.go
```

### Usando Docker Compose

```bash
# Iniciar os serviços
docker-compose up --build

# Parar os serviços
docker-compose down
```

### Usando Makefile

```bash
# Executar testes
make test

# Executar aplicação
make run

# Build Docker
make docker-build
make docker-run
```

## Endpoints

### GET /weather/{cep}

Retorna o clima atual para o CEP informado.

**Parâmetros:**
- `cep`: CEP brasileiro com 8 dígitos (pode incluir hífen)

**Respostas:**

| Status | Descrição | Exemplo |
|--------|-----------|---------|
| 200 | Sucesso | `{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.5}` |
| 422 | CEP inválido | `{"message": "invalid zipcode"}` |
| 404 | CEP não encontrado | `{"message": "can not find zipcode"}` |

**Exemplos de uso:**

```bash
# Consultar clima por CEP
curl http://localhost:8080/weather/01310100

# Consultar clima por CEP com hífen
curl http://localhost:8080/weather/01310-100
```

### GET /health

Health check endpoint.

```bash
curl http://localhost:8080/health
```

## Fórmulas de Conversão

- **Celsius para Fahrenheit:** `F = C * 1.8 + 32`
- **Celsius para Kelvin:** `K = C + 273`

## Testes

```bash
# Executar todos os testes
go test -v ./...

# Executar testes com cobertura
go test -v -cover ./...
```

## Deploy no Google Cloud Run

### Usando gcloud CLI

```bash
# Autenticar no Google Cloud
gcloud auth login

# Configurar projeto
gcloud config set project SEU_PROJECT_ID

# Deploy
gcloud run deploy weather-by-cep \
  --source . \
  --region us-central1 \
  --platform managed \
  --allow-unauthenticated \
  --set-env-vars "WEATHER_API_KEY=SUA_CHAVE"
```

### Usando Cloud Build

1. Configure a substituição `_WEATHER_API_KEY` no Cloud Build
2. Execute o build:

```bash
gcloud builds submit --config cloudbuild.yaml \
  --substitutions=_WEATHER_API_KEY=sua_chave
```

## Estrutura do Projeto

```
weather-by-cep/
├── main.go                     # Entrada principal
├── go.mod                      # Módulo Go
├── Dockerfile                  # Imagem Docker
├── docker-compose.yml          # Orquestração Docker
├── cloudbuild.yaml             # Config Cloud Build
├── Makefile                    # Comandos úteis
├── env.example                 # Exemplo de variáveis
├── README.md                   # Este arquivo
└── internal/
    ├── handlers/
    │   ├── weather.go          # Handler HTTP
    │   └── weather_test.go     # Testes do handler
    ├── models/
    │   └── models.go           # Modelos de dados
    └── services/
        ├── interfaces.go       # Interfaces para DI
        ├── viacep.go           # Serviço ViaCEP
        ├── viacep_test.go      # Testes ViaCEP
        ├── weather.go          # Serviço Weather
        └── weather_test.go     # Testes Weather
```

## APIs Utilizadas

- **ViaCEP:** https://viacep.com.br/ - Consulta de CEPs brasileiros
- **WeatherAPI:** https://www.weatherapi.com/ - Dados climáticos

## Licença

MIT
