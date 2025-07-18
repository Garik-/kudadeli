FROM node:24-alpine AS frontend

WORKDIR /app
COPY /frontend/package*.json ./
RUN npm ci
COPY ./frontend ./
RUN npm run build-only

FROM golang:1.24-alpine AS backend

WORKDIR /app

RUN cp /etc/ssl/certs/ca-certificates.crt /app/ca-certificates.crt

COPY ./bot/go.* ./
RUN go mod download

COPY ./bot ./
COPY --from=frontend /bot/web/public ./web/public

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server main.go

FROM scratch

WORKDIR /app

COPY --from=backend /app/server .
COPY --from=backend /app/ca-certificates.crt .

ENV SSL_CERT_FILE=/app/ca-certificates.crt

EXPOSE 8080

CMD ["./server"]