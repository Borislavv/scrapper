FROM --platform=linux/amd64 golang:1.22

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Устанавливаем необходимые зависимости и компиляторы
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

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENTRYPOINT ["go", "run", "./cmd/main.go"]
