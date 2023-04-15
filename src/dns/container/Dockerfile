FROM golang:latest
WORKDIR /app

COPY main.go .
RUN go build main.go

CMD ["./main"]