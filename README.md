# Flow Sight

[![License: Private](https://img.shields.io/badge/License-Private-red.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![Next.js Version](https://img.shields.io/badge/Next.js-15.4.4-black.svg)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue.svg)](https://www.typescriptlang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.5-blue.svg)](https://postgresql.org/)

**Flow Sight**は、個人の金融管理を行うためのモダンなWebアプリケーションです。クレジットカードや銀行口座の管理、収入・支出の追跡、キャッシュフロー予測など、健全な資金管理をサポートします。

## ✨ 主要機能

### 💰 金融データ管理
- **銀行口座管理** - 複数口座の残高管理
- **クレジットカード管理** - 締め日・支払日を考慮した管理
- **収入管理** - 月額固定・一時的収入の管理
- **固定支出管理** - 家賃、保険料などの定期支払い管理
- **月次利用額管理** - クレジットカードの月次総額管理

### 📊 予測・分析機能
- **キャッシュフロー予測** - 最大36ヶ月先までの資金残高推移予測
- **締め日・支払日考慮** - 正確な支払いスケジュール計算
- **日次残高推移** - 詳細な資金動向の可視化

## 🏗️ 技術スタック

### バックエンド
- **言語**: Go 1.24.5
- **フレームワーク**: Gin (HTTP Web Framework)
- **データベース**: PostgreSQL 17.5
- **ORM**: database/sql (標準ライブラリ)
- **マイグレーション**: golang-migrate
- **ドキュメント**: Swagger/OpenAPI
- **ログ**: slog (構造化ログ)

### フロントエンド
- **フレームワーク**: Next.js 15.4.4 (App Router)
- **言語**: TypeScript
- **UIライブラリ**: shadcn/ui (Radix UI + Tailwind CSS)
- **フォーム**: React Hook Form + Zod
- **通知**: Sonner
- **アイコン**: Lucide React
- **パッケージマネージャー**: pnpm

### インフラ・DevOps
- **コンテナ**: Docker & Docker Compose
- **オーケストレーション**: Kubernetes & Helm
- **プロキシ**: Nginx
- **環境管理**: .env ファイル

## 🚀 クイックスタート

### 前提条件
- Docker & Docker Compose
- Git
- (任意) Go 1.24.5, Node.js 18+, pnpm

### 1. プロジェクトクローン

```bash
git clone https://github.com/Soli0222/flow-sight.git
cd flow-sight
```

### 2. 環境設定

```bash
# 環境変数ファイルのセットアップ
make setup
```

`.env` を必要に応じて調整してください。

### 3. アプリケーション起動

```bash
# 全サービス（DB、バックエンド、フロントエンド、Nginx）を起動
make up

# 開発モード（ログ表示あり）
make dev
```

### 4. アクセス

- **アプリケーション**: http://localhost:4000
- **API ドキュメント**: http://localhost:4000/swagger/index.html
- **バックエンド API**: http://localhost:4000/api/v1
- **フロントエンド**: http://localhost:4000 (Nginx経由)

## 📁 プロジェクト構成

```
flow-sight/
├── backend/                    # Go APIサーバー
│   ├── cmd/                    # アプリケーションエントリーポイント
│   ├── internal/               # 内部パッケージ
│   │   ├── api/               # APIサーバー設定
│   │   ├── handlers/          # HTTPハンドラー
│   │   ├── models/            # データモデル
│   │   ├── repositories/      # データリポジトリ層
│   │   └── services/          # ビジネスロジック層
│   ├── migrations/            # データベースマイグレーション
│   └── docs/                  # Swagger生成ファイル
├── frontend/                  # Next.js フロントエンド
│   ├── src/                   # ソースコード
│   │   ├── app/              # App Router
│   │   ├── components/       # Reactコンポーネント
│   │   ├── lib/              # ユーティリティ
│   │   └── types/            # TypeScript型定義
│   └── public/               # 静的ファイル
├── helm-chart/               # Kubernetes Helmチャート
│   ├── templates/            # Kubernetesマニフェスト
│   └── values.yaml          # Helm設定値
├── documents/               # プロジェクトドキュメント
│   ├── api_specification.md # API仕様書
│   └── fromtend.md         # フロントエンド仕様書
├── docker-compose.yml      # Docker Compose設定
├── nginx.conf              # Nginx設定
└── Makefile               # 開発用コマンド
```

## 🛠️ 開発コマンド

### 基本操作

```bash
# ヘルプ表示
make help

# アプリケーション起動
make up

# アプリケーション停止
make down

# ログ表示
make logs

# 全サービス再起動
make restart

# 環境クリーンアップ
make clean
```

### サービス別ログ

```bash
# バックエンドログ
make logs-backend

# フロントエンドログ
make logs-frontend

# Nginxログ
make logs-nginx
```

### 開発用シェル

```bash
# バックエンドコンテナにアクセス
make backend-shell

# フロントエンドコンテナにアクセス
make frontend-shell
```

## 🧪 テスト

### バックエンドテスト

```bash
cd backend
go test -v ./...

# カバレッジレポート付き
go test -v -cover ./...
```

### フロントエンドテスト

```bash
cd frontend
pnpm test

# ウォッチモード
pnpm test:watch

# カバレッジレポート
pnpm test:coverage
```

## 📖 API ドキュメント

アプリケーション起動後、Swagger UIで詳細なAPI仕様を確認できます：

- **Swagger UI**: http://localhost:4000/swagger/index.html

### 主要APIエンドポイント

- **銀行口座**: `/api/v1/bank-accounts`
- **クレジットカード**: `/api/v1/credit-cards`
- **収入管理**: `/api/v1/income-sources`, `/api/v1/monthly-income-records`
- **固定支出**: `/api/v1/recurring-payments`
- **キャッシュフロー予測**: `/api/v1/cashflow-projection`
- **アプリ設定**: `/api/v1/settings`

## 🚢 デプロイメント

### Kubernetesデプロイメント

```bash
# Helmチャートでデプロイ
cd helm-chart
helm install flow-sight .

# アップグレード
helm upgrade flow-sight .

# アンインストール
helm uninstall flow-sight
```

### 環境変数設定

本番環境では以下の環境変数を適切に設定してください：

```bash
ENV=production

# データベース
DB_HOST=本番データベースホスト
DB_PASSWORD=強力なパスワード
DB_SSLMODE=require
```

## 🎯 パフォーマンス

### 要件
- **一般的なCRUD操作**: 500ms以内
- **キャッシュフロー予測計算**: 2秒以内
- **36ヶ月分の予測**: 高速計算対応

### 特徴
- 構造化ログによる監視
- ヘルスチェックエンドポイント
- レスポンシブデザイン対応

## 🔒 セキュリティ

- HTTPS推奨: 本番環境での暗号化通信
- 環境変数管理: 機密情報の適切な管理
- 構造化ログ: 重要イベントの記録

## 🗺️ 今後の拡張予定

1. **支出詳細管理** - 個別支出記録の管理
2. **予算管理機能** - 月次・年次予算設定
3. **レポート機能** - 分析とトレンド表示
4. **外部連携** - 銀行API、カード明細自動取得
5. **通知システム** - 支払い期限アラート
6. **モバイルアプリ** - React Native対応

## 📄 ライセンス

このプロジェクトはプライベートプロジェクトです。

## 🤝 貢献

プルリクエストやイシューの報告は歓迎します。以下のガイドラインに従ってください：

1. フォークしてブランチを作成
2. 変更を実装しテストを追加
3. コードスタイルガイドラインに準拠
4. プルリクエストを作成

## 📞 サポート

質問や問題がある場合は、GitHubのIssueページでお知らせください。

---

**Flow Sight** - 健全な資金管理をサポートする個人金融管理アプリケーション