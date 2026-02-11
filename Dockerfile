FROM golang:1.25.5-alpine AS builder

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . .

RUN go build -o blue-gopher ./cmd

FROM golang:1.25.5 AS naive-runner

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download

COPY . .

RUN go build -o blue-gopher ./cmd

CMD [ "./blue-gopher" ]

FROM alpine:latest AS runner 

WORKDIR /api

RUN mkdir -p internal/database

COPY --from=builder /build/internal/database/*.sql /api/internal/database/

COPY --from=builder /build/blue-gopher .

CMD [ "./blue-gopher" ]