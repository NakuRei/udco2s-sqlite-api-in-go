# ビルドステージ
FROM golang:1.22-alpine3.19 as builder

WORKDIR /app

# Goの依存関係をコピー＆インストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# CGOを有効にしてGoアプリケーションをビルド
# SQLiteの開発ライブラリはビルドに必要なためインストールする
# build-baseとsqlite-devはビルド後には不要なため、このステージでのみインストール
RUN apk add --no-cache build-base sqlite-dev \
  && CGO_ENABLED=1 GOOS=linux go build -o udco2s-sqlite-api-in-go \
  && apk del build-base sqlite-dev

# 実行ステージ
FROM alpine:3.19

ARG USERNAME=user
ARG ENV=prod

# NOTE: SQLiteのランタイムライブラリをインストール
RUN apk add --no-cache sqlite-libs

RUN addgroup -S $USERNAME && adduser -S -G $USERNAME -h /home/$USERNAME $USERNAME

# RUN if [ "$ENV" = "dev" ]; then \
#   addgroup $USERNAME && adduser -D -G $USERNAME -h /home/$USERNAME $USERNAME; \
#   else \
#   addgroup -S $USERNAME && adduser -S -G $USERNAME -h /home/$USERNAME $USERNAME; \
#   fi

WORKDIR /home/$USERNAME/workspace

# ビルドステージからバイナリファイルをコピー
COPY --from=builder --chown=$USERNAME:$USERNAME /app/udco2s-sqlite-api-in-go .

USER $USERNAME

EXPOSE 8080

CMD ["./udco2s-sqlite-api-in-go"]
