# 🚀 クイックスタートガイド

炎上シミュレーターを5分で起動する手順です。

## 📋 前提条件

- ✅ Docker Desktop がインストール・起動済み
- ✅ Git がインストール済み
- ✅ gcloud CLI がインストール済み（[インストール手順](https://cloud.google.com/sdk/docs/install)）
- ⚠️ GCP プロジェクトと認証設定（後述）

## 🏃 最速起動手順

### 1. リポジトリをクローン

```bash
# GitHubからクローン
git clone https://github.com/Tattsum/enjo.git
cd enjo
```

### 2. GCP認証情報を設定

Application Default Credentials (ADC) を取得します：

```bash
# GCPにログイン
gcloud auth login

# Application Default Credentialsを取得
gcloud auth application-default login

# 認証情報ファイルをプロジェクトにコピー
cp ~/.config/gcloud/application_default_credentials.json backend/
```

### 3. 環境変数ファイルを設定

`backend/.env` を編集してGCPプロジェクトIDを設定：

```bash
# エディタで開く
nano backend/.env
```

以下を確認・編集：

```env
GCP_PROJECT_ID=your-project-id  # ← あなたのGCPプロジェクトIDに変更
GCP_LOCATION=us-central1
PORT=8080
```

フロントエンドの環境変数も設定（既に存在する場合はスキップ）：

```bash
cp frontend/.env.local.example frontend/.env.local
```

### 4. Vertex AI APIを有効化

```bash
# あなたのプロジェクトIDを設定
export PROJECT_ID="your-project-id"

# Vertex AI APIを有効化
gcloud services enable aiplatform.googleapis.com --project=$PROJECT_ID
```

### 5. Docker Composeで起動

```bash
# すべてのサービスを起動（初回は5-10分かかります）
docker-compose up --build
```

または、バックグラウンドで起動：

```bash
docker-compose up --build -d
```

### 6. アクセス

ブラウザで以下のURLを開いてください：

- **フロントエンド**: http://localhost:3000
- **GraphQL Playground**: http://localhost:8080/graphql
- **バックエンド ヘルスチェック**: http://localhost:8080/health

---

## 🧪 動作確認

### 1. ヘルスチェック

```bash
curl http://localhost:8080/health
```

期待される出力：
```json
{"status":"OK"}
```

### 2. GraphQL クエリ

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"{ health }"}'
```

期待される出力：
```json
{"data":{"health":"OK"}}
```

### 3. フロントエンド確認

ブラウザで http://localhost:3000 を開いて：

1. テキストエリアに「今日はいい天気ですね」と入力
2. 炎上レベルを「3」に設定
3. 「🔥 炎上化する」ボタンをクリック
4. 変換結果が表示されるか確認

---

## 🛑 停止・クリーンアップ

### サービスを停止

```bash
# 停止（コンテナは保持）
docker-compose stop

# 停止して削除
docker-compose down

# ボリュームも含めて完全削除
docker-compose down -v
```

### ログ確認

```bash
# すべてのログ
docker-compose logs -f

# バックエンドのみ
docker-compose logs -f backend

# フロントエンドのみ
docker-compose logs -f frontend
```

---

## 🐛 トラブルシューティング

### ポートが使用中

```bash
# 既に8080ポートが使われている場合
lsof -i :8080

# または3000ポート
lsof -i :3000

# プロセスを停止してから再起動
```

### Dockerビルドエラー

```bash
# キャッシュをクリアして再ビルド
docker-compose build --no-cache

# イメージを削除して再作成
docker-compose down --rmi all
docker-compose up --build
```

### 認証エラー

ブラウザのコンソールまたはバックエンドログで以下を確認：

```bash
docker-compose logs backend | grep -i "vertex\|gcp\|auth"
```

エラー例：

- `GCP_PROJECT_ID is required` → .envファイルのプロジェクトIDを確認
- `Failed to create Vertex AI client` → ADC認証情報を確認
- `application_default_credentials.json: no such file` → 手順2の認証情報コピーを実行
- `Permission denied` → GCPプロジェクトでVertex AI APIが有効か確認

### ADC認証情報の再取得

認証エラーが続く場合：

```bash
# 再度ADCを取得
gcloud auth application-default login

# 認証情報を再コピー
cp ~/.config/gcloud/application_default_credentials.json backend/

# Dockerサービスを再起動
docker-compose restart backend
```

---

## 📚 次のステップ

### 開発者向け

- [PROJECT_RULES.md](../PROJECT_RULES.md) - TDD開発ルール
- [README.md](../README.md) - 詳細なドキュメント

### テストを実行

```bash
# すべてのテスト
make test

# バックエンドのみ
make backend-test

# フロントエンドのみ
make frontend-test
```

### コードチェック

```bash
# すべてのチェック（フォーマット、Lint、テスト）
make check

# バックエンドのみ
make backend-check

# フロントエンドのみ
make frontend-check
```

---

## 💡 ヒント

- **初回起動は時間がかかります**: 依存関係のダウンロードに5-10分かかることがあります
- **ホットリロード対応**: コードを変更すると自動的に再読み込みされます
- **認証情報の管理**: `application_default_credentials.json` は機密情報です。Gitにコミットしないでください（.gitignoreで除外済み）
- **Vertex AI料金**: 無料枠がありますが、使用量に応じて課金される場合があります。[料金ページ](https://cloud.google.com/vertex-ai/pricing)を確認してください
- **教育目的での使用**: 実際のSNSでの悪用は厳禁です

---

## 🆘 サポート

問題が発生した場合：

1. [Issues](https://github.com/Tattsum/enjo/issues) で既知の問題を検索
2. 新しいIssueを作成して質問
3. ログを確認: `docker-compose logs -f`

---

**Have Fun! 🔥**
