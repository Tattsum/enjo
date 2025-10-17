# ローカル環境セットアップガイド

このガイドでは、画像生成機能を含む炎上シミュレーターをローカル環境で実行する方法を説明します。

## 前提条件

- Docker & Docker Compose がインストール済み
- GCPプロジェクトへのアクセス権限
- （オプション）Twitter API認証情報

## セットアップ手順

### 1. GCP認証情報の設定

#### 方法1: Application Default Credentials (ADC) を使用（推奨）

```bash
# Google Cloud CLIにログイン
gcloud auth application-default login

# 認証情報ファイルを確認
ls ~/.config/gcloud/application_default_credentials.json

# プロジェクトのbackendディレクトリにコピー
cp ~/.config/gcloud/application_default_credentials.json \
   backend/application_default_credentials.json
```

#### 方法2: サービスアカウントキーを使用

```bash
# GCPコンソールでサービスアカウントキーを作成してダウンロード
# ファイルを backend/application_default_credentials.json として保存
```

### 2. 環境変数の設定

```bash
# backend/.env.example をコピー
cp backend/.env.example backend/.env

# エディタで backend/.env を編集
```

**backend/.env の設定例:**

```env
# GCP Configuration for Vertex AI
GCP_PROJECT_ID=your-actual-project-id  # 実際のプロジェクトIDに変更
GCP_LOCATION=us-central1

# Server Port
PORT=8080

# Twitter API Configuration (オプション - 投稿機能を使う場合のみ)
TWITTER_API_KEY=your_twitter_api_key_here
TWITTER_API_SECRET=your_twitter_api_secret_here
TWITTER_ACCESS_TOKEN=your_twitter_access_token_here
TWITTER_ACCESS_TOKEN_SECRET=your_twitter_access_token_secret_here
```

### 3. GCP APIの有効化

```bash
# Vertex AI APIを有効化
gcloud services enable aiplatform.googleapis.com

# 確認
gcloud services list --enabled | grep aiplatform
```

### 4. Docker環境の起動

```bash
# プロジェクトルートディレクトリで実行
docker-compose up --build
```

初回起動時は以下の処理が実行されます：
- バックエンドの依存パッケージのインストール
- フロントエンドの依存パッケージのインストール
- 両方のサービスの起動

**起動完了のログ:**
```
frontend_1  | ✓ Ready in 2.5s
backend_1   | Server is running on http://localhost:8080
```

### 5. アプリケーションへのアクセス

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **GraphQL Playground**: http://localhost:8080/graphql

## 機能の確認

### 1. 炎上テキスト生成

1. http://localhost:3000 にアクセス
2. 普通のテキストを入力（例：「今日のランチは美味しかった」）
3. 炎上度スライダーを調整（1-5）
4. 「🔥 炎上化する」ボタンをクリック
5. 炎上テキストが生成される

### 2. 画像生成（新機能！）

1. 炎上テキスト生成後、「🎨 画像を生成」セクションが表示される
2. スタイルを選択（ミーム風/リアル調/イラスト調/ドラマチック）
3. 「🎨 画像を生成」ボタンをクリック
4. 5-15秒程度で画像が生成される
5. 生成された画像をダウンロードまたは再生成可能

### 3. Twitter投稿（オプション）

1. Twitter API認証情報が設定されている場合
2. 「𝕏 Xに投稿」または「𝕏 画像付きで投稿」ボタンをクリック
3. Twitterに投稿される

## トラブルシューティング

### 1. 認証エラー

**エラー:**
```
Error: google: could not find default credentials
```

**解決方法:**
```bash
# 認証情報ファイルが正しい場所にあるか確認
ls backend/application_default_credentials.json

# ない場合は再度コピー
cp ~/.config/gcloud/application_default_credentials.json \
   backend/application_default_credentials.json

# Docker再起動
docker-compose restart backend
```

### 2. Vertex AI APIエラー

**エラー:**
```
Error: Vertex AI API has not been used in project
```

**解決方法:**
```bash
# APIを有効化
gcloud services enable aiplatform.googleapis.com

# プロジェクトIDを確認
gcloud config get-value project

# .envファイルのGCP_PROJECT_IDと一致するか確認
cat backend/.env | grep GCP_PROJECT_ID
```

### 3. 画像生成が遅い

**原因:**
- 初回リクエストはコールドスタート（10-20秒）
- ネットワーク状況により5-15秒かかる

**対策:**
- 小さい画像サイズ（512x512）を使用
- 2回目以降は高速化される

### 4. ポート競合エラー

**エラー:**
```
Error: port 3000 is already in use
```

**解決方法:**
```bash
# 使用中のプロセスを確認
lsof -i :3000
lsof -i :8080

# プロセスを停止してから再起動
docker-compose down
docker-compose up
```

### 5. Docker コンテナが起動しない

**解決方法:**
```bash
# ログを確認
docker-compose logs backend
docker-compose logs frontend

# コンテナを完全にクリーンアップして再起動
docker-compose down -v
docker-compose up --build
```

## 開発モード

### バックエンドのホットリロード

バックエンドのコードを変更すると、自動的に再起動されます。

```bash
# ログを確認
docker-compose logs -f backend
```

### フロントエンドのホットリロード

フロントエンドのコードを変更すると、自動的にリロードされます。

```bash
# ログを確認
docker-compose logs -f frontend
```

### テストの実行

```bash
# バックエンドテスト
cd backend
make backend-check

# フロントエンドテスト
cd frontend
npm run lint
npm run type-check
npm run test
```

## コスト管理

### 画像生成のコスト

ローカル環境でも実際のVertex AI APIを使用するため、コストが発生します。

| 操作 | コスト（512x512） | コスト（1024x1024） |
|------|------------------|-------------------|
| 1回の画像生成 | $0.020 | $0.040 |
| 10回の画像生成 | $0.20 | $0.40 |
| 100回の画像生成 | $2.00 | $4.00 |

### コスト削減のヒント

1. **開発中は512x512を使用**: 品質は十分、コストは半分
2. **不要な再生成を避ける**: 1回の生成で満足できる結果を得る
3. **統合テストは慎重に実行**: `RUN_INTEGRATION_TESTS=true`を設定しない限り実行されない
4. **GCPコンソールで使用量を監視**:
   - https://console.cloud.google.com/billing
   - Vertex AI使用量を定期的に確認

## 環境のシャットダウン

```bash
# コンテナを停止（データは保持）
docker-compose stop

# コンテナを停止して削除
docker-compose down

# ボリュームも含めて完全に削除
docker-compose down -v
```

## 次のステップ

- [画像生成機能の詳細](docs/FEATURE_IMAGE_GENERATION.md)
- [Twitter投稿機能の詳細](docs/FEATURE_TWITTER_POST.md)
- [プロジェクトルール](PROJECT_RULES.md)
- [統合テストの実行](backend/INTEGRATION_TEST_README.md)

## サポート

問題が発生した場合：

1. このREADMEのトラブルシューティングセクションを確認
2. Docker・GCPのログを確認
3. 環境変数が正しく設定されているか確認
4. GitHubのIssuesで質問

---

**注意:** この環境は開発・テスト用です。本番環境へのデプロイには別途設定が必要です。
