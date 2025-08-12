FROM golang:1.24.2-bookworm

WORKDIR /src

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        libvips-dev pkg-config && \
    rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /main .

EXPOSE 7000

ENTRYPOINT ["/main"]