FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o /app/users
FROM scratch
COPY --from=builder /app/users /app/users
ENTRYPOINT [ "/app/users" ]