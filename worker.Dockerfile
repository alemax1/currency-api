FROM golang:1.22

WORKDIR /currencies-worker

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .

RUN go build -o /app ./cmd/currencies-worker/main.go

CMD ["/app", "-c", ".env"]