FROM golang:1.22-alpine as build


WORKDIR /app

COPY . /app

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./haotrader cmd/main/main.go



FROM alpine:latest

USER root

COPY --from=build /app/haotrader /app/haotrader
RUN chmod +x /app/haotrader

EXPOSE 8080

WORKDIR /app

CMD ["./haotrader"]