-- データベース初期化用SQLファイル
-- PostgreSQL用の初期設定

-- 拡張機能の有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- タイムゾーン設定
SET timezone = 'Asia/Tokyo';

-- 基本的なテーブル作成（後でマイグレーションで詳細を作成）
-- このファイルは基本的な設定のみ、実際のテーブルはマイグレーションで管理

-- ユーザー用のスキーマ作成
CREATE SCHEMA IF NOT EXISTS digeon;

-- インデックス用の設定
-- 日本語検索用の設定
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- 全文検索用の設定
CREATE EXTENSION IF NOT EXISTS unaccent;

-- 権限設定
GRANT ALL PRIVILEGES ON DATABASE digeon_db TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA digeon TO postgres;

-- 統計情報の有効化
ALTER DATABASE digeon_db SET log_statement = 'all';
ALTER DATABASE digeon_db SET log_min_duration_statement = 1000;