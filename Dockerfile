FROM golang:1.24 AS builder

ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -buildid=" -o /zen-lang ./main.go

FROM scratch AS final

LABEL org.opencontainers.image.source="https://github.com/Senither/zen-lang"
LABEL org.opencontainers.image.title="zen-lang"
LABEL org.opencontainers.image.description="Zen language interpreter binary"

COPY --from=builder /zen-lang /zen-lang

ENTRYPOINT ["/zen-lang"]
