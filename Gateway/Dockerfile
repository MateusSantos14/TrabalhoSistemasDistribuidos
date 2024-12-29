# Use uma imagem Go oficial
FROM golang:1.21

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie os arquivos go.mod e go.sum para o contêiner
COPY go.mod go.sum ./

# Baixe as dependências do Go
RUN go mod download

# Copie o restante do código para o contêiner
COPY . .

# Gere os arquivos .go a partir do .proto usando o protoc
RUN apt-get update && apt-get install -y protobuf-compiler

# Compile o aplicativo Go
RUN go build -o gateway ./gateway.go

# Exponha a porta do contêiner (se necessário)
EXPOSE 8888

# Comando para rodar o aplicativo
CMD ["./gateway"]
