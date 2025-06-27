FROM golang:1.24-alpine

WORKDIR /app

COPY . .

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

RUN go mod tidy

RUN go build -o cart cmd/main.go

CMD [ "./cart" ]
