FROM golang:1.23-alpine

WORKDIR /app

# 依存関係をコピーしてインストール
COPY go.mod ./
RUN go mod download

# ソースコードをコピー
COPY . .

# .air.tomlファイルの作成（既存のコード）
RUN if [ ! -f .air.toml ]; then \
    echo 'root = "."' > .air.toml && \
    echo 'tmp_dir = "tmp"' >> .air.toml && \
    echo '[build]' >> .air.toml && \
    # ... other echo commands ... \
    echo '  clear_on_rebuild = true' >> .air.toml \
    ; fi

# この行を追加
RUN go mod tidy

# アプリをビルド
RUN go build -o main .

# 開発モードでホットリロードを有効にするためのair（Go用の開発ツール）をインストール
RUN go install github.com/air-verse/air@latest

EXPOSE 8080

# 開発環境では air を使用
CMD ["air", "-c", ".air.toml"]
