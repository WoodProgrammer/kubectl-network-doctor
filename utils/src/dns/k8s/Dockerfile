FROM golang:1.20.4-bullseye as builder
WORKDIR /app

COPY cmd/main.go .
RUN go build main.go 
RUN ls /app/

FROM debian:stable-slim
COPY --from=builder /app/main /opt/
RUN ls
WORKDIR /opt/
CMD ["./main"]