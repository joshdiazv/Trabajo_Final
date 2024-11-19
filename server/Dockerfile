FROM golang:1.20-alpine

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos del servidor
COPY server.go . 

# Copia los archivos CSV desde la carpeta server
COPY movies.csv ./movies.csv
COPY ratings.csv ./ratings.csv

# Compila el servidor
RUN go build -o server server.go

# Exponer el puerto del servidor
EXPOSE 8080

# Ejecuta el servidor
CMD ["./server"]
