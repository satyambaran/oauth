FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o /app/clients
FROM scratch
COPY --from=builder /app/clients /app/clients
ENTRYPOINT [ "/app/clients" ]