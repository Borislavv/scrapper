FROM --platform=linux/amd64 golang:1.22 as builder

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    curl \
    gnupg \
    gcc \
    g++ \
    libc6-dev \
    linux-headers-generic \
    make \
    g++-x86-64-linux-gnu \
    libc6-dev-amd64-cross \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/go-spider

COPY . .

RUN go mod download

RUN --mount=type=cache,target=/root/.cache/go-build,id=golang_skeleton_toolkit_go_mod_cache,uid=0,gid=0  \
    --mount=type=secret,id=netrc,target=/root/.netrc \
    go mod download

RUN --mount=type=secret,id=netrc,target=/root/.netrc \
    go build -o spider ./cmd/main.go

RUN #go build -o spider ./cmd/main.go

ENTRYPOINT ["/go/src/go-spider/spider"]