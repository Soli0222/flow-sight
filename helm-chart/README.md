# Flow-Sight Helm Chart

Flow-Sightアプリケーション用のHelmチャートです。

## 概要

このHelmチャートは以下のコンポーネントをデプロイします：

- **Database**: PostgreSQL 17.5
- **Backend**: Go製のAPIサーバー（ポート8080）
- **Frontend**: Next.jsフロントエンド（ポート3000）
- **Ingress**: パスベースルーティング
  - `/api/*` → Backend
  - `/swagger/*` → Backend  
  - `/*` → Frontend

## 機能

### データベース起動待機
バックエンドのDeploymentには、データベースの起動を待機するinitContainerが含まれています：

- `pg_isready`コマンドでデータベースの準備完了を確認
- データベースが利用可能になるまでバックエンドの起動を待機
- 内部データベース・外部データベース両方に対応
- `backend.initContainer.enabled: false`で無効化可能

## 前提条件

- Kubernetes 1.16+
- Helm 3.0+
- Ingress Controller（nginx-ingress、Traefik等）

## インストール方法

1. リポジトリをクローン
```bash
git clone <repository-url>
cd flow-sight/helm-chart
```

2. 依存関係の更新
```bash
helm dependency update
```

3. values.yamlの設定
```bash
# バックエンドのデータベース接続設定
backend:
  environment:
    DB_HOST: "your-database-host"  # 外部DB使用時はホスト名を変更
  database:
    name: flowsight_db
    user: postgres
    port: 5432
  secrets:
    externalName: ""  # 外部Secretを使う場合はSecret名を指定
    DB_PASSWORD: "your-secure-database-password"

# 内部データベース（開発用、本番では通常無効）
database:
  enabled: true  # 外部DB使用時は false

# Ingress設定
ingress:
  hosts:
    - host: your-domain.com
      paths:
        - path: /api
          pathType: Prefix
          backend:
            service: backend
        - path: /swagger
          pathType: Prefix
          backend:
            service: backend
        - path: /
          pathType: Prefix
          backend:
            service: frontend
```

4. インストール
```bash
helm install flow-sight .
```

## 設定オプション

### Database設定
- `database.enabled`: 内部データベースの有効/無効（開発用、本番では通常false）
- `database.persistence.enabled`: 永続化の有効/無効
- `database.persistence.size`: ストレージサイズ

### Backend設定
- `backend.replicaCount`: レプリカ数
- `backend.image.repository`: イメージリポジトリ
- `backend.environment.DB_HOST`: データベースホスト（外部DB使用時に変更）
- `backend.database.name`: データベース名
- `backend.database.user`: データベースユーザー
- `backend.database.port`: データベースポート
- `backend.initContainer.enabled`: データベース待機用initContainerの有効/無効
- `backend.initContainer.image.*`: initContainerで使用するPostgreSQLイメージ設定
- `backend.secrets.externalName`: 外部Secretの名前（設定すると外部Secretを参照、内部Secret作成をスキップ）
- `backend.secrets.DB_PASSWORD`: データベースパスワード（Secretで管理）

### Frontend設定
- `frontend.replicaCount`: レプリカ数
- `frontend.image.repository`: イメージリポジトリ

### Ingress設定
- `ingress.enabled`: Ingressの有効/無効
- `ingress.className`: IngressClass名
- `ingress.hosts`: ホスト設定

## 使用例

### 開発環境（内部データベース使用）
```bash
# values.yaml
database:
  enabled: true  # 内部PostgreSQLを使用
backend:
  environment:
    DB_HOST: "flow-sight-db"  # 内部DBサービス名
  database:
    name: flowsight_db
    user: postgres
    port: 5432
  initContainer:
    enabled: true  # DB起動待機を有効
  secrets:
    DB_PASSWORD: "development-password"
```

### 本番環境（外部データベース使用）
```bash
# values-production.yaml
database:
  enabled: false  # 内部DBは無効
backend:
  environment:
    DB_HOST: "prod-postgres.example.com"  # 外部DBホスト
  database:
    name: flowsight_production
    user: flowsight_user
    port: 5432
  initContainer:
    enabled: true  # 外部DBでも起動待機を有効
  secrets:
    DB_PASSWORD: "secure-production-password"
```

### 外部Secret使用例（本番環境推奨）
```bash
# 事前に外部Secretを作成
kubectl create secret generic flow-sight-secrets \
  --from-literal=DB_PASSWORD="secure-production-password"

# values.yaml
backend:
  secrets:
    externalName: "flow-sight-secrets"  # 外部Secretを参照
    # 他のsecret値は外部Secretで管理されるため不要
```

## アップグレード

```bash
helm upgrade flow-sight .
```

## アンインストール

```bash
helm uninstall flow-sight
```

## トラブルシューティング

### ポッド状態の確認
```bash
kubectl get pods
kubectl describe pod <pod-name>
```

### ログの確認
```bash
kubectl logs <pod-name>
```

### サービスの確認
```bash
kubectl get svc
```

### Ingressの確認
```bash
kubectl get ingress
kubectl describe ingress flow-sight
```
