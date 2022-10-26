FROM golang:alpine

LABEL mainteiner="customers"

WORKDIR /customers

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main src/cmd/main.go src/cmd/build_api.go

RUN find . -name "*.go" -type f -delete

CMD ["./main"]
