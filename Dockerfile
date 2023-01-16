# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o /run-app /app/cmd/run-server/main.go

EXPOSE 8080

CMD [ "/run-app", "-addr", "0.0.0.0:8080" ]