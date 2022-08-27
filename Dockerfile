#build stage building the app here
FROM golang:1.18.5 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o bank-app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bank-app .
COPY app.env .

EXPOSE 8080
CMD ["/app/bank-app"]
