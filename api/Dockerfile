FROM golang:1.20-alpine

WORKDIR /app

COPY api.go . 

RUN go build -o api api.go

CMD ["./api"]
