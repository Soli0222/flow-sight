.PHONY: help up down logs build clean setup dev

# デフォルトのターゲット
help: ## このヘルプメッセージを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## 初期セットアップ（.envファイルをコピー）
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "📝 .envファイルを作成しました。Google OAuth情報を設定してください。"; \
	else \
		echo "✅ .envファイルは既に存在します。"; \
	fi

up: ## アプリケーション全体を起動
	docker-compose up -d

build: ## イメージをビルドして起動
	docker-compose up -d --build

dev: ## 開発モードで起動（ログ表示あり）
	docker-compose up --build

down: ## アプリケーション全体を停止
	docker-compose down

logs: ## 全サービスのログを表示
	docker-compose logs -f

logs-backend: ## バックエンドのログを表示
	docker-compose logs -f backend

logs-frontend: ## フロントエンドのログを表示
	docker-compose logs -f frontend

logs-nginx: ## Nginxのログを表示
	docker-compose logs -f nginx

clean: ## 全てのコンテナ、イメージ、ボリュームを削除
	docker-compose down -v --rmi all

restart: ## アプリケーションを再起動
	docker-compose restart

# 個別サービスの管理
backend-shell: ## バックエンドコンテナにシェル接続
	docker-compose exec backend sh

frontend-shell: ## フロントエンドコンテナにシェル接続
	docker-compose exec frontend sh

db-shell: ## データベースに接続
	docker-compose exec db psql -U postgres -d flowsight_db
