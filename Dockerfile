# multi-stage building
FROM golang:1.25.0 AS build-stage 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/server

# final stage / production image
FROM scratch AS production-stage
WORKDIR /app

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-stage /app/server .

CMD ["./server"]

