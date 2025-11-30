FROM golang:1.24 AS builder

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN set -eux; \
    LDFLAGS="-s -w \
      -X github.com/senither/zen-lang/cli.Version=$VERSION \
      -X github.com/senither/zen-lang/cli.GitCommit=$GIT_COMMIT \
      -X github.com/senither/zen-lang/cli.BuildDate=$BUILD_DATE"; \
    GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$LDFLAGS" -o /zen-lang ./main.go

FROM scratch AS final

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_DATE=unknown

LABEL org.opencontainers.image.source="https://github.com/Senither/zen-lang"
LABEL org.opencontainers.image.title="zen-lang"
LABEL org.opencontainers.image.description="Zen language interpreter binary"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.revision="$GIT_COMMIT"
LABEL org.opencontainers.image.created="$BUILD_DATE"

COPY --from=builder /zen-lang /zen-lang

ENTRYPOINT ["/zen-lang"]
