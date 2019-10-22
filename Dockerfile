# Эта секция для сборки бинарника
FROM golang:1.12 AS compile

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

# Эта секция будет рабочей в дальнейшем
FROM alpine:latest

# tini - корневой процесс с pid = 1, он отвечает за нормальное завершение дочерних процессов, в нашем случае -
# приложения
RUN apk add tini

ENTRYPOINT ["tini", "--"]

WORKDIR /app
COPY --from=compile /app/server /app/server
COPY --from=compile /app/configs /app/configs
COPY --from=compile /app/migrations /app/migrations

EXPOSE 8084

CMD ["/app/server"]

