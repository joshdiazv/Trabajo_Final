FROM golang:1.20-alpine

WORKDIR /app

# Copia el código fuente
COPY http_proxy.go .

# Compila el archivo Go
RUN go build -o proxy http_proxy.go

# Comando para ejecutar el binario compilado
CMD ["./proxy"]
