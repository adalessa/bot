FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

ADD . .

RUN go build -o /bot

CMD ["/bot"]