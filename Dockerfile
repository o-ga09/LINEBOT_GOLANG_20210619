#デプロイ用コンテナに含めるバイナリを作成するコンテナ
FROM golang:1.18.2-bullseye as deploy-builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

#-----------------------------------------------
#デプロイ用コンアテナ
FROM centos:centos7 as deploy

RUN yum -y update

EXPOSE "8080"

COPY --from=deploy-builder /app/app .

CMD ["./app"]