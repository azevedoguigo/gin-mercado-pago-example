FROM golang:1.22-alpine

WORKDIR /gin-mercado-pago-example

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

CMD [ "go", "run", "main.go" ]

EXPOSE 8080