FROM golang:1.16-alpine

ENV GO111MODULE=on

ENV GOPATH=/

COPY ./ ./

EXPOSE 8080

RUN apk --no-cache add curl openssl

RUN chmod 777 gen_ca.sh && chmod 777 gen_cert.sh

RUN ./gen_ca.sh

RUN go mod download

RUN go mod tidy

RUN go build -o kproxy ./cmd/main.go

CMD ./kproxy
