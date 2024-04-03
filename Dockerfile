FROM golang:1.21 AS compile
WORKDIR /app

COPY . .
RUN export GOPROXY=https://goproxy.cn && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o starland-account ./cmd/...


FROM ubuntu:20.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y ca-certificates openssh-client tzdata libimage-exiftool-perl && echo "Asia/Shanghai" > /etc/timezone && rm -f /etc/localtime && dpkg-reconfigure -f noninteractive tzdata

WORKDIR /app
COPY --from=compile /app/starland-account .

ENV PATH="/app:$PATH"

CMD ["starland-account"]
