# Digeon Frontend

X(旧Twitter)クローンアプリケーションのフロントエンド部分。React + TypeScriptで構築されています。

## 技術スタック

- **React** 18.x
- **TypeScript** 
- **React Router** - ルーティング
- **Axios** - HTTP通信
- **Create React App** - プロジェクト構成

## プロジェクト構造

```
src/
├── components/     # 再利用可能なコンポーネント
│   ├── Header.tsx
│   ├── PostCard.tsx
│   └── PostForm.tsx
├── pages/         # ページコンポーネント
│   ├── Home.tsx
│   ├── Login.tsx
│   └── Register.tsx
├── types/         # TypeScript型定義
│   └── index.ts
├── api/           # API通信
│   └── client.ts
├── hooks/         # カスタムフック（今後追加予定）
└── utils/         # ユーティリティ関数（今後追加予定）
```

## 主要機能

### 認証
- ユーザー登録
- ログイン/ログアウト
- JWT認証

### 投稿機能
- テキスト投稿（最大280文字）
- 画像投稿（最大4枚）
- 投稿削除
- リアルタイムタイムライン表示

### インタラクション
- いいね機能
- リポスト機能
- コメント機能（UI準備済み）

## 開発環境のセットアップ

### 前提条件
- Node.js 16.x以上
- npm または yarn

### インストール

```bash
# 依存関係のインストール
npm install

# 開発サーバーの起動
npm start
```

### 環境変数

`.env`ファイルを作成し、以下を設定：

```
REACT_APP_API_URL=http://localhost:8080/api
```

## 利用可能なスクリプト

### `npm start`
開発サーバーを起動します。  
ブラウザで [http://localhost:3000](http://localhost:3000) を開いてアプリケーションを確認できます。

### `npm test`
テストランナーを起動します。

### `npm run build`
本番用のアプリケーションをビルドします。

### `npm run eject`
Create React Appの設定を取り出します（非推奨）。

## API エンドポイント

### 認証
- `POST /api/auth/login` - ログイン
- `POST /api/auth/register` - ユーザー登録
- `POST /api/auth/logout` - ログアウト
- `GET /api/auth/me` - 現在のユーザー情報取得

### 投稿
- `GET /api/posts` - 投稿一覧取得
- `POST /api/posts` - 投稿作成
- `DELETE /api/posts/:id` - 投稿削除
- `POST /api/posts/:id/like` - いいね
- `DELETE /api/posts/:id/like` - いいね解除
- `POST /api/posts/:id/repost` - リポスト
- `DELETE /api/posts/:id/repost` - リポスト解除

### ユーザー
- `GET /api/users/:id` - ユーザー情報取得
- `GET /api/users/:id/posts` - ユーザーの投稿一覧
- `POST /api/users/:id/follow` - フォロー
- `DELETE /api/users/:id/follow` - フォロー解除

## コンポーネント設計

### Header
- ナビゲーションバー
- ログイン状態に応じた表示切り替え
- ログアウト機能

### PostCard
- 投稿の表示
- いいね・リポスト・コメントボタン
- 投稿者の情報表示
- 投稿削除機能（自分の投稿のみ）

### PostForm
- 新規投稿フォーム
- 文字数制限（280文字）
- 画像アップロード機能
- プレビュー機能

## 状態管理

現在はReactのuseStateを使用した簡単な状態管理を実装しています。  
将来的にはReduxやZustandなどの状態管理ライブラリの導入を検討しています。

## スタイリング

インラインスタイルを使用したシンプルなスタイリング。  
Twitterライクなデザインを目指しています。

## 今後の実装予定

### Phase 1（基本機能）
- [x] ユーザー認証
- [x] 投稿機能  
- [x] タイムライン表示
- [x] いいね機能
- [x] リポスト機能

### Phase 2（拡張機能）
- [ ] コメント機能の実装
- [ ] フォロー機能
- [ ] ユーザープロフィールページ
- [ ] 検索機能
- [ ] 通知機能

### Phase 3（最適化）
- [ ] 無限スクロール
- [ ] リアルタイム更新（WebSocket）
- [ ] PWA対応
- [ ] テスト追加
- [ ] パフォーマンス最適化

## 開発ガイドライン

### コーディング規約
- TypeScriptの型定義を活用
- コンポーネントは関数コンポーネントで作成
- propsの型定義は必須
- ファイル名はPascalCaseで統一

### Git運用
- feature/機能名でブランチを作成
- コミットメッセージは日本語で記述
- Pull Request前にビルドエラーがないことを確認

## トラブルシューティング

### よくある問題

1. **API接続エラー**
   - バックエンドサーバーが起動しているか確認
   - .envファイルのAPI URLが正しいか確認

2. **ビルドエラー**
   - node_modulesを削除してnpm installを実行
   - TypeScriptの型エラーを確認

3. **認証エラー**
   - ローカルストレージのtokenをクリア
   - ブラウザの開発者ツールでネットワークタブを確認

## ライセンス

MIT License