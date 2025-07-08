# Digeon Backend

X(旧Twitter)クローンアプリのバックエンド API

## 技術スタック

- **言語**: Go 1.24
- **フレームワーク**: Echo v4
- **データベース**: PostgreSQL
- **認証**: JWT
- **アーキテクチャ**: Clean Architecture

## プロジェクト構成

```
.
├── cmd/
│   └── server/           # アプリケーションエントリポイント
├── internal/
│   ├── domain/           # ドメイン層
│   │   ├── entity/       # エンティティ
│   │   ├── repository/   # リポジトリインターフェース
│   │   └── service/      # ドメインサービス
│   ├── application/      # アプリケーション層
│   │   ├── usecase/      # ユースケース
│   │   └── dto/          # データ転送オブジェクト
│   ├── infrastructure/   # インフラストラクチャ層
│   │   ├── database/     # データベース実装
│   │   └── web/          # 外部API連携
│   └── presentation/     # プレゼンテーション層
│       ├── handler/      # HTTPハンドラー
│       ├── middleware/   # ミドルウェア
│       └── router/       # ルーティング
├── config/               # 設定ファイル
├── migrations/           # データベースマイグレーション
└── docs/                 # ドキュメント
```

## 主要機能

### Phase 1: 基本機能
- [x] 環境構築
- [ ] ユーザー認証（登録、ログイン、ログアウト）
- [ ] 投稿機能（作成、取得、削除）
- [ ] タイムライン表示

### Phase 2: インタラクション
- [ ] いいね機能
- [ ] フォロー機能
- [ ] コメント機能

### Phase 3: 仕上げ
- [ ] 検索機能
- [ ] 通知機能
- [ ] UI/UX調整

## API エンドポイント

### 認証
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/login` - ログイン
- `POST /api/auth/logout` - ログアウト
- `POST /api/auth/refresh` - トークンリフレッシュ

### 投稿
- `GET /api/posts` - 投稿一覧取得
- `POST /api/posts` - 投稿作成
- `PUT /api/posts/:id` - 投稿更新
- `DELETE /api/posts/:id` - 投稿削除

### ユーザー
- `GET /api/users/:id` - ユーザー情報取得
- `PUT /api/users/:id` - ユーザー情報更新
- `GET /api/users/:id/posts` - ユーザーの投稿一覧
- `GET /api/users/:id/followers` - フォロワー一覧
- `GET /api/users/:id/following` - フォロー中一覧

### インタラクション
- `POST /api/posts/:id/like` - いいね
- `DELETE /api/posts/:id/like` - いいね解除
- `POST /api/posts/:id/repost` - リポスト
- `DELETE /api/posts/:id/repost` - リポスト解除

### フォロー
- `POST /api/users/:id/follow` - フォロー
- `DELETE /api/users/:id/follow` - フォロー解除

## セットアップ

### 前提条件
- Go 1.24以上
- PostgreSQL 14以上

### 開発環境構築

1. リポジトリクローン
```bash
git clone <repository-url>
cd feat-back
```

2. 依存関係インストール
```bash
go mod tidy
```

3. 環境変数設定
```bash
cp .env.example .env
# .envファイルを編集
```

4. データベース設定
```bash
# PostgreSQL起動
# データベース作成
createdb digeon_db
```

5. マイグレーション実行
```bash
go run cmd/migrate/main.go up
```

6. アプリケーション起動
```bash
go run cmd/server/main.go
```

## テスト

```bash
go test ./...
```

## 開発時間割

- **午前中**: 基本機能（認証、投稿、タイムライン）
- **午後前半**: インタラクション（いいね、フォロー、コメント）
- **午後後半**: 仕上げ（検索、通知、調整）

## 制約事項

- 投稿文字数制限: 280文字
- 画像サイズ制限: 5MB
- 開発期間: 1日