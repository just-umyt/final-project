FROM golang:1.24-alpine

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV GOPROXY=https://mirrors.aliyun.com/goproxy/

RUN go mod tidy


RUN go build -tags musl -ldflags '-extldflags "-static"' -o metrics cmd/main.go

CMD [ "./metrics" ]