FROM golang:1.24

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN ./scripts/generate_ssl_cert.sh
RUN go build -v -o /usr/local/bin/server ./cmd/server

CMD ["server"]
