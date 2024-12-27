# Utilizando a imagem oficial do Golang para construir o projeto
FROM golang:1.21-alpine

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos go.mod e go.sum e baixar as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código da aplicação
COPY . .

# Construir o binário
RUN go build -o /complete-api

# Expõe a porta que a aplicação irá usar
EXPOSE 4000

# Comando para rodar a aplicação
CMD ["/complete-api"]
