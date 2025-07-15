FROM golang:1.24-alpine AS builder

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV GOPROXY=https://goproxy.cn

RUN go mod tidy

RUN go build -tags musl -ldflags '-extldflags "-static"' -o metrics cmd/main.go

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/metrics .

CMD [ "./metrics" ]
