FROM golang:1.24.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache gcc musl-dev

RUN go build -o /godocker

EXPOSE 8080

CMD ["/godocker"]