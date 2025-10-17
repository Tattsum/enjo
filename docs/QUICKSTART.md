# 🚀 クイックスタートガイド

炎上シミュレーターを5分で起動する手順です。

## 📋 前提条件

- ✅ Docker Desktop がインストール・起動済み
- ✅ Git がインストール済み
- ⚠️ Gemini API キー（後述）

## 🏃 最速起動手順（API キーなし）

### 1. リポジトリをクローン

```bash
# GitHubからクローン
git clone https://github.com/Tattsum/enjo.git
cd enjo
```

### 2. 環境変数ファイルを作成

```bash
# Makefileを使って自動作成
make setup
```

または手動で：

```bash
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
```

### 3. Docker Composeで起動

```bash
# すべてのサービスを起動（初回は5-10分かかります）
docker-compose up --build
```

または、バックグラウンドで起動：

```bash
docker-compose up --build -d
```

### 4. アクセス

ブラウザで以下のURLを開いてください：

- **フロントエンド**: http://localhost:3000
- **GraphQL Playground**: http://localhost:8080/graphql
- **バックエンド ヘルスチェック**: http://localhost:8080/health

## ⚠️ 注意事項

**API キーなしの場合:**
- UIは表示されますが、実際のテキスト変換機能は動作しません
- エラーメッセージ「GEMINI_API_KEY is not set」が表示されます

**完全な機能を使うには:**
- 次のセクション「Gemini API キーの設定」を参照してください

---

## 🔑 Gemini API キーの設定（完全版）

実際にテキスト変換機能を使用するには、Gemini API キーが必要です。

### 方法1: Google AI Studio で取得（最も簡単）

1. [Google AI Studio](https://aistudio.google.com/apikey) にアクセス
2. 「Get API Key」または「Create API Key」をクリック
3. APIキーをコピー

### 方法2: gcloud コマンドで取得

詳細は [docs/SETUP_API_KEY.md](./SETUP_API_KEY.md) を参照してください。

### API キーを設定

`backend/.env` ファイルを編集：

```bash
# エディタで開く
nano backend/.env

# または
vim backend/.env
```

以下のように修正：

```env
# 取得したAPIキーに置き換える
GEMINI_API_KEY=AIzaSy...your_actual_api_key_here
PORT=8080
```

### サービスを再起動

```bash
# バックエンドのみ再起動
docker-compose restart backend

# すべてのログを確認
docker-compose logs -f
```

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

### API キーエラー

ブラウザのコンソールまたはバックエンドログで以下を確認：

```bash
docker-compose logs backend | grep -i "gemini\|api"
```

エラー例：
- `GEMINI_API_KEY is not set` → .envファイルを確認
- `Invalid API key` → APIキーが正しいか確認
- `API quota exceeded` → 使用制限を超えています

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
- **API制限に注意**: Gemini APIには無料枠の制限があります（1分あたり15リクエスト）
- **教育目的での使用**: 実際のSNSでの悪用は厳禁です

---

## 🆘 サポート

問題が発生した場合：

1. [Issues](https://github.com/Tattsum/enjo/issues) で既知の問題を検索
2. 新しいIssueを作成して質問
3. ログを確認: `docker-compose logs -f`

---

**Have Fun! 🔥**
