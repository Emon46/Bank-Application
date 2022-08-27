#build stage building the app here
FROM golang:1.18.5 AS builder
WORKDIR /app
RUN apt update && apt install curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
COPY . .
RUN CGO_ENABLED=0 go build -o bank-app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bank-app .
COPY --from=builder /app/migrate .
COPY app.env .
COPY ./db/migration ./migration
COPY ./start.sh .
COPY ./wait-for.sh .
EXPOSE 8080
CMD ["/app/bank-app"]
ENTRYPOINT ["/app/start.sh"]
