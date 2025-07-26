# Flow Sight

Flow Sightは個人の金融管理を行うためのWebアプリケーションです。クレジットカードやローンの支払い管理、月次キャッシュフロー予測、収入・支出の管理を通じて、健全な資金管理をサポートします。

## 技術スタック

### フロントエンド
- **フレームワーク**: Next.js 15.4 (React 19)
- **言語**: TypeScript
- **スタイル**: Tailwind CSS
- **UI**: Radix UI + shadcn/ui
- **状態管理**: React Hooks + Context API
- **認証**: Google OAuth 2.0 + JWT

### バックエンド
- **言語**: Go 1.21
- **フレームワーク**: Gin
- **データベース**: PostgreSQL 15
- **認証**: JWT + Google OAuth 2.0
- **API仕様**: OpenAPI/Swagger

### インフラ
- **コンテナ**: Docker & Docker Compose
- **リバースプロキシ**: Nginx
- **統合ポート**: localhost:4000

## クイックスタート

### 前提条件
- Docker & Docker Compose
- Google Cloud Console アカウント

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd flow-sight
```

### 2. Google OAuth設定

1. [Google Cloud Console](https://console.cloud.google.com/)でプロジェクトを作成
2. OAuth 2.0認証情報を作成:
   - **承認済みのJavaScript生成元**: `http://localhost:4000`
   - **承認済みのリダイレクトURI**: `http://localhost:4000/api/v1/auth/google/callback`
3. クライアントIDとクライアントシークレットを取得

### 3. 環境変数の設定

```bash
# 環境変数ファイルの作成
make setup

# .envファイルを編集してGoogle OAuth情報を設定
GOOGLE_CLIENT_ID=your-google-client-id-here
GOOGLE_CLIENT_SECRET=your-google-client-secret-here
```

### 4. アプリケーションの起動

```bash
# 本番と同じ構成で起動（推奨）
make build

# または通常起動
make up

# 開発時のログ表示
make dev
```

### 5. アクセス

ブラウザで http://localhost:4000 にアクセスしてください。

- **フロントエンド**: http://localhost:4000
- **API**: http://localhost:4000/api/v1
- **API ドキュメント**: http://localhost:4000/swagger/index.html

## Docker構成の特徴

- **統一された構成**: 開発・本番で同じDockerfileとdocker-compose.ymlを使用
- **マルチステージビルド**: 最適化されたフロントエンドビルド
- **Node.js Standalone**: Next.js 15の最適化されたスタンドアロン出力
- **セキュリティ**: 非rootユーザーでの実行
- **効率的なレイヤーキャッシュ**: 依存関係とソースコードの分離

## 主要機能

### 金融管理機能
1. **銀行口座管理** - 複数の銀行口座の残高管理
2. **資産管理** - クレジットカードやローンの管理
3. **収入管理** - 給与や副収入の管理
4. **定期支払い管理** - 家賃、保険料などの固定費管理
5. **カード利用額管理** - 月次クレジットカード利用額の記録
6. **キャッシュフロー予測** - 最大36ヶ月先までの資金推移予測

### システム機能
- **Google OAuth認証** - 簡単ログイン
- **レスポンシブデザイン** - PC・モバイル対応
- **ダークモード** - テーマ切り替え
- **リアルタイム同期** - データの自動更新

## 開発用コマンド

```bash
# アプリケーション管理
make up          # 起動
make down        # 停止
make restart     # 再起動
make build       # ビルドして起動
make dev         # 開発時（ログ表示あり）
make clean       # 全て削除
make clean       # 全て削除

# ログ確認
make logs        # 全サービス
make logs-backend   # バックエンド
make logs-frontend  # フロントエンド
make logs-nginx     # Nginx

# デバッグ
make backend-shell  # バックエンドコンテナに接続
make frontend-shell # フロントエンドコンテナに接続
make db-shell       # データベースに接続
```

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────┐
│                    localhost:4000                          │
├─────────────────────────────────────────────────────────────┤
│                     Nginx                                   │
│  ┌─────────────────┬─────────────────────────────────────┐  │
│  │   Frontend      │          Backend API                │  │
│  │  (Next.js)      │           (Go/Gin)                  │  │
│  │     :3000       │            :8080                    │  │
│  └─────────────────┴─────────────────────────────────────┘  │
│                           │                                 │
│                    ┌─────────────┐                         │
│                    │ PostgreSQL  │                         │
│                    │    :5432    │                         │
│                    └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

### ルーティング設定

- `/` → フロントエンド (Next.js)
- `/api/*` → バックエンド API (Go)
- `/swagger/*` → API ドキュメント

## プロジェクト構成

```
flow-sight/
├── frontend/              # Next.js フロントエンド
│   ├── src/
│   │   ├── app/          # App Router
│   │   ├── components/   # React コンポーネント
│   │   ├── lib/         # ユーティリティ
│   │   └── types/       # TypeScript 型定義
│   └── Dockerfile.dev   # 開発用 Dockerfile
│
├── backend/              # Go バックエンド
│   ├── cmd/             # エントリーポイント
│   ├── internal/        # アプリケーションコード
│   ├── migrations/      # DB マイグレーション
│   └── Dockerfile       # Dockerfile
│
├── docker-compose.yml   # Docker Compose 設定
├── nginx.conf          # Nginx 設定
├── Makefile           # 開発用コマンド
└── README.md          # このファイル
```

## Google OAuth設定詳細

### 承認済みドメイン設定

Google Cloud Consoleで以下の設定を行ってください：

1. **OAuth同意画面**:
   - アプリケーション名: Flow Sight
   - 承認済みドメイン: localhost

2. **OAuth 2.0 クライアントID**:
   - アプリケーションの種類: ウェブアプリケーション
   - 承認済みのJavaScript生成元:
     - `http://localhost:4000`
   - 承認済みのリダイレクトURI:
     - `http://localhost:4000/api/v1/auth/google/callback`

### 本番環境での設定

本番環境では以下のドメインを設定してください：
- JavaScript生成元: `https://yourdomain.com`
- リダイレクトURI: `https://yourdomain.com/api/v1/auth/google/callback`

## トラブルシューティング

### よくある問題

**1. OAuth認証エラー**
```
error: redirect_uri_mismatch
```
→ Google Cloud ConsoleのリダイレクトURIを確認してください

**2. API接続エラー**
```
fetch failed
```
→ バックエンドサービスが起動していることを確認してください

**3. データベース接続エラー**
```
connection refused
```
→ データベースサービスが起動していることを確認してください

### ログの確認

```bash
# 全サービスのログ
make logs

# 特定のサービス
make logs-backend
make logs-frontend
make logs-nginx
```

### データベースのリセット

```bash
# データベースを含めて全て削除
make clean

# 再起動
make build
```

## ライセンス

このプロジェクトはプライベートプロジェクトです。

## サポート

質問や問題がある場合は、プロジェクトのIssueページでお知らせください。
