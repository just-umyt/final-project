FROM golang:1.24-alpine

RUN apk add --no-cache \
    librdkafka-dev \
    pkgconf \
    build-base

WORKDIR /app

COPY . .

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

RUN go mod tidy


RUN go build -o metrics cmd/main.go

CMD [ "./metrics" ]

# FROM golang:1.24-alpine AS builder


# WORKDIR /app
# COPY . .
# RUN go mod download

# RUN GOOS=linux GOARCH=amd64 go build -tags musl -o /app/bin/main .

# FROM ubuntu

# RUN  sudo apt install -y librdkafka-dev pkg-config build-essential

# WORKDIR /app

# COPY --from=builder /app/bin/main /app/bin/main

# ENTRYPOINT ["/app/bin/main"]
