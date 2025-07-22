FROM golang:1.24-alpine

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

# ENV GOPROXY=https://goproxy.cn

# ENV GOPROXY=https://proxy.golang.org,direct


RUN go mod tidy

RUN go build -tags musl -ldflags '-extldflags "-static"' -o cart cmd/main.go

CMD [ "./cart",  "-env=prod"]