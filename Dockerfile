FROM golang:1.22.1

WORKDIR /app

COPY . /app

RUN go mod tidy
RUN go build -o main main.go

CMD ["./main"]
