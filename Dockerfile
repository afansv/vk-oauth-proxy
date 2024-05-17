FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/vkoauthproxy/main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]
