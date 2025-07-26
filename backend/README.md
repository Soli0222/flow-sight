# Flow Sight Backend API

Flow Sightは個人の金融管理を行うためのWebアプリケーションのバックエンドAPIです。クレジットカードやローンの支払い管理、月次キャッシュフロー予測、収入・支出の管理を通じて、健全な資金管理をサポートします。

## 技術スタック

- **言語**: Go 1.21
- **フレームワーク**: Gin (HTTP Web Framework)
- **データベース**: PostgreSQL 15
- **ORM**: database/sql (標準ライブラリ)
- **マイグレーション**: golang-migrate
- **ドキュメント**: Swagger/OpenAPI
- **コンテナ**: Docker & Docker Compose

## API機能

### 主要機能

1. **資産管理API** - クレジットカードやローンの管理
2. **銀行口座管理API** - 銀行口座と残高の管理
3. **収入管理API** - 収入源と月次収入実績の管理
4. **固定支出管理API** - 家賃、保険料などの定期支払いの管理
5. **カード月次利用額管理API** - クレジットカードの月次利用総額の管理
6. **キャッシュフロー予測API** - 将来の資金残高推移の予測計算
7. **アプリケーション設定API** - ユーザー設定の管理

### キャッシュフロー予測の特徴

- 締め日・支払日を考慮した正確な支払いスケジュール計算
- 複数銀行口座の残高統合
- 最大36ヶ月先までの予測
- 日次残高推移の詳細計算

## セットアップ

### 前提条件

- Go 1.21以上
- Docker & Docker Compose
- PostgreSQL 15（ローカル開発の場合）

### 環境変数

`.env`ファイルを作成して以下の環境変数を設定してください：

```bash
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=flowsight_db
DB_SSLMODE=disable
JWT_SECRET=your-jwt-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
ENV=development
```

### Google OAuth設定

1. [Google Cloud Console](https://console.cloud.google.com/)でプロジェクトを作成
2. OAuth 2.0認証情報を作成:
   - 承認済みのリダイレクトURI: `http://localhost:4000/api/v1/auth/google/callback`
   - 承認済みのJavaScript生成元: `http://localhost:4000`
3. クライアントIDとクライアントシークレットを環境変数に設定

### Docker統合環境での起動

プロジェクトルートから：

```bash
# 初期セットアップ
make setup

# 統合環境起動（フロント・バック・DB全て）
make up

# アクセス: http://localhost:4000
```

詳細な統合環境の使い方は、プロジェクトルートの`README.md`を参照してください。
DB_PASSWORD=password
DB_NAME=flowsight_db
DB_SSLMODE=disable
JWT_SECRET=your-jwt-secret-key
ENV=development
```

### Docker Composeを使用した起動

```bash
# 依存関係のインストール
make deps

# Docker環境での起動
make dev-setup
```

### ローカル開発

```bash
# 依存関係のインストール
go mod download

# データベースの起動
docker-compose up -d db

# マイグレーションの実行
make migrate-up

# アプリケーションの起動
make run
```

## API エンドポイント

### クレジットカード管理
- `GET /api/v1/credit-cards` - クレジットカード一覧取得
- `POST /api/v1/credit-cards` - クレジットカード登録
- `GET /api/v1/credit-cards/{id}` - クレジットカード詳細取得
- `PUT /api/v1/credit-cards/{id}` - クレジットカード更新
- `DELETE /api/v1/credit-cards/{id}` - クレジットカード削除

### 銀行口座管理
- `GET /api/v1/bank-accounts` - 口座一覧取得
- `POST /api/v1/bank-accounts` - 口座登録
- `GET /api/v1/bank-accounts/{id}` - 口座詳細取得
- `PUT /api/v1/bank-accounts/{id}` - 口座更新
- `DELETE /api/v1/bank-accounts/{id}` - 口座削除

### 収入管理
- `GET /api/v1/income-sources` - 収入源一覧取得
- `POST /api/v1/income-sources` - 収入源登録
- `GET /api/v1/monthly-income-records` - 月次収入実績一覧取得
- `POST /api/v1/monthly-income-records` - 月次収入実績登録

### 固定支出管理
- `GET /api/v1/recurring-payments` - 固定支出一覧取得
- `POST /api/v1/recurring-payments` - 固定支出登録
- `PUT /api/v1/recurring-payments/{id}` - 固定支出更新
- `DELETE /api/v1/recurring-payments/{id}` - 固定支出削除

### カード月次利用額管理
- `GET /api/v1/card-monthly-totals` - カード月次利用額一覧取得
- `POST /api/v1/card-monthly-totals` - カード月次利用額登録
- `PUT /api/v1/card-monthly-totals/{id}` - カード月次利用額更新
- `DELETE /api/v1/card-monthly-totals/{id}` - カード月次利用額削除

### キャッシュフロー予測
- `GET /api/v1/cashflow-projection` - キャッシュフロー予測取得

### アプリケーション設定
- `GET /api/v1/settings` - 設定一覧取得
- `PUT /api/v1/settings` - 設定更新

### その他
- `GET /api/v1/health` - ヘルスチェック
- `GET /swagger/*` - API ドキュメント

## 開発

### ディレクトリ構造

```
backend/
├── cmd/                    # アプリケーションエントリーポイント
│   └── main.go
├── internal/               # 内部パッケージ
│   ├── api/               # APIサーバー設定
│   ├── config/            # 設定管理
│   ├── database/          # データベース接続
│   ├── handlers/          # HTTPハンドラー
│   ├── models/            # データモデル
│   ├── repositories/      # データリポジトリ層
│   └── services/          # ビジネスロジック層
├── migrations/            # データベースマイグレーション
├── docker-compose.yml     # Docker Compose設定
├── Dockerfile            # Docker設定
├── go.mod               # Go モジュール
├── go.sum               # Go モジュールチェックサム
├── Makefile             # 開発用コマンド
└── README.md            # このファイル
```

### 開発用コマンド

```bash
# アプリケーションの起動
make run

# テストの実行
make test

# コードフォーマット
make fmt

# Lintの実行
make lint

# Swaggerドキュメント生成
make swagger

# Docker環境での起動
make docker-run

# Docker環境の停止
make docker-down

# マイグレーションの実行
make migrate-up

# マイグレーションのロールバック
make migrate-down
```

### テスト

```bash
# すべてのテストを実行
go test -v ./...

# カバレッジレポート付きでテスト実行
go test -v -cover ./...
```

### APIドキュメント

アプリケーション起動後、以下のURLでSwagger UIにアクセスできます：

```
http://localhost:8080/swagger/index.html
```

## 本番環境

### 環境変数

本番環境では以下の環境変数を適切に設定してください：

- `ENV=production`
- `JWT_SECRET` - 強力なランダム文字列
- `DB_*` - 本番データベースの接続情報

### セキュリティ

- JWTシークレットは十分にランダムで複雑なものを使用
- データベース接続情報は環境変数で管理
- HTTPS通信の使用を推奨
- 定期的なセキュリティアップデート

### パフォーマンス要件

- 一般的なCRUD操作：500ms以内
- キャッシュフロー予測計算：2秒以内
- 36ヶ月分のキャッシュフロー予測を高速計算

## ライセンス

このプロジェクトはプライベートプロジェクトです。

## 貢献

プルリクエストやイシューの報告は歓迎します。コードスタイルガイドラインに従って開発してください。

## サポート

質問や問題がある場合は、プロジェクトのIssueページでお知らせください。
