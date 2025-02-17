# プロジェクト計画

## 必須機能

### 1. 記事管理

- 記事の作成、編集、削除機能
- Markdownサポートとプレビュー機能
- オプション: Headless CMS（ContentfulやStrapiなど）の使用を検討すると、コンテンツ管理が容易になる

### 2. コメント機能

- 記事へのコメント投稿機能
- コメントの管理（削除、通報など）

### 3. カテゴリーとタグ

- 記事のカテゴリー分け
- タグ機能

### 4. 検索機能

- 記事のタイトルや内容から検索する機能

### 5. レスポンシブデザイン

- モバイル、タブレット、デスクトップに対応

## 追加機能

### 1. RSSフィード

- 記事のRSSフィード配信

### 2. アナリティクス

- Google Analyticsなどの統合

### 3. SEO対策

- メタタグ、タイトルタグ、ディスクリプションなどの設定管理

### 4. パフォーマンス最適化

- 画像の遅延読み込みやキャッシュ戦略

### 5. インテグレーション

- GitHub Gistsの統合によるコードスニペットの埋め込み機能

### 6. 多言語サポート

- UIや記事の多言語対応

## 実装ポイント

### フロントエンド

- ReactとTypeScriptで構築し、型の安全性を維持

### バックエンド

- Go言語でAPIを構築、高速で効率良いサーバー実現
- RESTful APIを設計

## データベーステーブルの設計

### テーブル構造


```sql
users テーブル
CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
);
articles テーブル
CREATE TABLE articles (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  author_id INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  published BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (author_id) REFERENCES users(id)
);
comments テーブル
CREATE TABLE comments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  article_id INT,
  user_id INT,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  approved BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (article_id) REFERENCES articles(id),
  FOREIGN KEY (user_id) REFERENCES users(id)
);
categories テーブル
CREATE TABLE categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
tags テーブル
CREATE TABLE tags (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
article_categories テーブル
CREATE TABLE article_categories (
  id INT AUTO_INCREMENT PRIMARY KEY,
  article_id INT,
  category_id INT,
  FOREIGN KEY (article_id) REFERENCES articles(id),
  FOREIGN KEY (category_id) REFERENCES categories(id)
);
article_tags テーブル
CREATE TABLE article_tags (
  id INT AUTO_INCREMENT PRIMARY KEY,
  article_id INT,
  tag_id INT,
  FOREIGN KEY (article_id) REFERENCES articles(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);
analytics テーブル
CREATE TABLE analytics (
  id INT AUTO_INCREMENT PRIMARY KEY,
  article_id INT,
  views INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (article_id) REFERENCES articles(id)
);
translations テーブル
CREATE TABLE translations (
  id INT AUTO_INCREMENT PRIMARY KEY,
  article_id INT,
  lang VARCHAR(10),
  translated_title VARCHAR(255),
  translated_content TEXT,
  FOREIGN KEY (article_id) REFERENCES articles(id)
);