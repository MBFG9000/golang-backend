FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o myapp ./cmd/api/main.go
RUN go build -o seeder ./cmd/seeder/main.go

FROM alpine
COPY --from=builder /app/myapp /myapp
COPY --from=builder /app/seeder /seeder

ENTRYPOINT ["/myapp"]
