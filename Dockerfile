FROM golang:1.22

WORKDIR /currencies-api

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o /app ./cmd/currencies-api/main.go

CMD ["/app", "-c", ".env"]
