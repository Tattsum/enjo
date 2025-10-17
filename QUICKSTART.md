# クイックスタートガイド

最速でローカル環境を起動する手順です。

## ✅ 前提条件チェック

すでに以下が設定済みです：
- ✅ `backend/.env` ファイルが存在
- ✅ `backend/application_default_credentials.json` が存在
- ✅ Docker & Docker Compose がインストール済み

## 🚀 起動手順（3ステップ）

### 1. 環境変数の確認

```bash
# GCP_PROJECT_ID が正しく設定されているか確認
cat backend/.env | grep GCP_PROJECT_ID
```

**期待される出力:**
```
GCP_PROJECT_ID=your-actual-project-id
```

もし `your_gcp_project_id_here` のままなら、実際のプロジェクトIDに変更してください：

```bash
# プロジェクトIDを確認
gcloud config get-value project

# .env を編集
vi backend/.env  # または好きなエディタ
```

### 2. Docker起動

```bash
# プロジェクトルートで実行
docker-compose up --build
```

**初回起動は5-10分程度かかります**（依存パッケージのインストール）

起動完了のログ:
```
frontend_1  | ✓ Ready in 2.5s
frontend_1  | - Local:        http://localhost:3000
backend_1   | Server is running on http://localhost:8080
```

### 3. ブラウザでアクセス

http://localhost:3000 を開く

## 🎨 画像生成機能を試す

1. **テキスト入力**
   ```
   今日のランチは最高でした！
   ```

2. **炎上度を選択**（スライダー: 3）

3. **「🔥 炎上化する」をクリック**
   - 炎上テキストが生成される

4. **「🎨 画像を生成」をクリック**（新機能！）
   - スタイルを選択（ミーム風/リアル調など）
   - 「🎨 画像を生成」ボタンをクリック
   - **5-15秒で画像が生成されます**

5. **生成された画像を確認**
   - ダウンロード可能
   - 再生成も可能
   - Twitter投稿も可能（認証情報設定済みの場合）

## 💰 コスト注意

画像生成は実際のVertex AI APIを使用します：
- 1回の生成: 約$0.02（512x512）
- 10回の生成: 約$0.20
- 開発中は小さいサイズを推奨

## 🛑 停止方法

```bash
# Ctrl+C で停止、または別ターミナルで
docker-compose down
```

## 🔧 トラブルシューティング

### 画像生成エラーが出る場合

1. **GCP_PROJECT_IDを確認**
   ```bash
   cat backend/.env | grep GCP_PROJECT_ID
   ```

2. **Vertex AI APIが有効か確認**
   ```bash
   gcloud services list --enabled | grep aiplatform
   ```

   有効でない場合:
   ```bash
   gcloud services enable aiplatform.googleapis.com
   ```

3. **Dockerを再起動**
   ```bash
   docker-compose restart backend
   ```

### ポート競合エラー

```bash
# 使用中のプロセスを確認
lsof -i :3000
lsof -i :8080

# Docker再起動
docker-compose down
docker-compose up
```

### 初回起動が遅い

- 正常です！依存パッケージのインストールに時間がかかります
- 2回目以降は高速です

## 📚 詳細ドキュメント

- [詳細なセットアップガイド](LOCAL_SETUP.md)
- [画像生成機能の詳細](docs/FEATURE_IMAGE_GENERATION.md)
- [統合テストの実行方法](backend/INTEGRATION_TEST_README.md)

## ✨ 次のステップ

画像生成が成功したら：
1. 異なるスタイルを試す
2. 炎上度を変えて再生成
3. コードを確認（TDD実装済み！）
4. テストを実行して理解を深める

---

**問題が発生した場合:** [LOCAL_SETUP.md](LOCAL_SETUP.md) のトラブルシューティングセクションを参照
