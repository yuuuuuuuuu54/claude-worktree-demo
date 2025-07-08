# Makefile for Digeon Backend

# 変数定義
APP_NAME := digeon-backend
DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_DEV := docker-compose -f docker-compose.yml -f docker-compose.dev.yml

# デフォルトターゲット
.PHONY: help
help: ## このヘルプを表示
	@echo "使用可能なコマンド:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# 開発環境
.PHONY: dev
dev: ## 開発環境を起動
	$(DOCKER_COMPOSE_DEV) up --build

.PHONY: dev-down
dev-down: ## 開発環境を停止
	$(DOCKER_COMPOSE_DEV) down

.PHONY: dev-logs
dev-logs: ## 開発環境のログを表示
	$(DOCKER_COMPOSE_DEV) logs -f api

# 本番環境
.PHONY: up
up: ## 本番環境を起動
	$(DOCKER_COMPOSE) up -d --build

.PHONY: down
down: ## 本番環境を停止
	$(DOCKER_COMPOSE) down

.PHONY: restart
restart: down up ## 本番環境を再起動

.PHONY: logs
logs: ## 本番環境のログを表示
	$(DOCKER_COMPOSE) logs -f api

# データベース
.PHONY: db-up
db-up: ## データベースのみ起動
	$(DOCKER_COMPOSE) up -d postgres redis

.PHONY: db-migrate
db-migrate: ## マイグレーションを実行
	$(DOCKER_COMPOSE) exec api go run cmd/migrate/main.go up

.PHONY: db-rollback
db-rollback: ## マイグレーションをロールバック
	$(DOCKER_COMPOSE) exec api go run cmd/migrate/main.go down

.PHONY: db-shell
db-shell: ## データベースシェルに接続
	$(DOCKER_COMPOSE) exec postgres psql -U postgres -d digeon_db

# テスト
.PHONY: test
test: ## テストを実行
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## テストカバレッジを表示
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ビルド
.PHONY: build
build: ## アプリケーションをビルド
	go build -o bin/$(APP_NAME) cmd/server/main.go

.PHONY: build-linux
build-linux: ## Linux用にビルド
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(APP_NAME)-linux cmd/server/main.go

# Docker
.PHONY: docker-build
docker-build: ## Dockerイメージをビルド
	docker build -t $(APP_NAME) .

.PHONY: docker-run
docker-run: ## Dockerコンテナを実行
	docker run -p 8080:8080 $(APP_NAME)

# クリーンアップ
.PHONY: clean
clean: ## 一時ファイルを削除
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	docker system prune -f

.PHONY: clean-all
clean-all: clean ## すべてのDockerリソースを削除
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE_DEV) down -v
	docker system prune -a -f

# 依存関係
.PHONY: deps
deps: ## 依存関係をインストール
	go mod tidy
	go mod download

.PHONY: deps-upgrade
deps-upgrade: ## 依存関係をアップグレード
	go get -u ./...
	go mod tidy

# リント
.PHONY: lint
lint: ## リントを実行
	golangci-lint run

.PHONY: fmt
fmt: ## コードをフォーマット
	go fmt ./...

# サーバー起動（ローカル）
.PHONY: run
run: ## ローカルでサーバーを起動
	go run cmd/server/main.go

# デバッグ
.PHONY: debug
debug: ## デバッグモードで起動
	go run cmd/server/main.go --debug

# 設定ファイル
.PHONY: config
config: ## 設定ファイルを作成
	cp .env.example .env

# 初期セットアップ
.PHONY: setup
setup: config deps ## 初期セットアップ
	@echo "セットアップが完了しました。"
	@echo "開発環境を起動するには: make dev"