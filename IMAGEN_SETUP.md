# Imagen API セットアップガイド

画像生成機能を使用するための詳細なセットアップガイドです。

## ✅ 修正完了（2025-10-17）

Internal エラーを修正し、Imagen API をREST API経由で直接呼び出す実装に変更しました。

### 主な変更点

1. **REST API実装に移行**: genai SDKのサポート不足のため、Vertex AI REST APIを直接使用
2. **正しいモデル名**: `imagegeneration@002` (Imagen 2) を使用
3. **認証**: Application Default Credentials (ADC) を使用

## 🔧 必須セットアップ

### 1. Vertex AI APIの有効化

```bash
# APIを有効化
gcloud services enable aiplatform.googleapis.com --project=tmp-rnd-ai

# 確認
gcloud services list --enabled --project=tmp-rnd-ai | grep aiplatform
```

**期待される出力:**
```
aiplatform.googleapis.com      Vertex AI API
```

### 2. 認証情報の確認

```bash
# ADCが設定されているか確認
ls -la backend/application_default_credentials.json

# 権限を確認
cat backend/application_default_credentials.json | jq .quota_project_id
```

### 3. プロジェクトIDの確認

```bash
# .envファイルを確認
cat backend/.env | grep GCP_PROJECT_ID
```

**現在の設定:**
```
GCP_PROJECT_ID=tmp-rnd-ai
GCP_LOCATION=us-central1
```

## 🎨 画像生成機能の使い方

### ブラウザからテスト

1. **アクセス**: http://localhost:3000
2. **テキスト入力**: 「今日のランチは最高でした！」
3. **炎上化**: 「🔥 炎上化する」をクリック
4. **画像生成**: 「🎨 画像を生成」をクリック
5. **待機**: 5-15秒で画像が生成されます

### GraphQL Playgroundからテスト

http://localhost:8080/graphql にアクセスして以下を実行:

```graphql
mutation {
  generateImage(input: {
    text: "今日のランチは最高でした！"
    style: MEME
    aspectRatio: SQUARE
  }) {
    imageUrl
    prompt
    generatedAt
  }
}
```

**成功レスポンス例:**
```json
{
  "data": {
    "generateImage": {
      "imageUrl": "data:image/png;base64,iVBORw0KGgo...",
      "prompt": "A dramatic scene of...",
      "generatedAt": "2025-10-17T14:30:00Z"
    }
  }
}
```

## 🐛 トラブルシューティング

### エラー: 422 (Unprocessable Entity)

**原因**: モデル名が間違っているか、APIが有効化されていない

**解決方法**:
```bash
# APIを再度有効化
gcloud services enable aiplatform.googleapis.com --project=tmp-rnd-ai

# 利用可能なモデルを確認（参考）
# Imagenは通常のモデルリストには表示されないため、ドキュメントを参照
```

### エラー: Internal error encountered

**原因**: Imagen APIの内部エラー（通常は一時的）

**解決方法**:
1. 数秒待ってリトライ
2. プロンプトをシンプルにする
3. 別のリージョンを試す（.envで`GCP_LOCATION`を変更）

### エラー: Permission denied

**原因**: 権限不足

**解決方法**:
```bash
# 現在のアカウントを確認
gcloud auth list

# プロジェクトのIAMロールを確認
gcloud projects get-iam-policy tmp-rnd-ai

# 必要に応じて権限を追加
gcloud projects add-iam-policy-binding tmp-rnd-ai \
  --member="user:YOUR_EMAIL" \
  --role="roles/aiplatform.user"
```

### エラー: Failed to get credentials

**原因**: 認証情報ファイルが見つからない

**解決方法**:
```bash
# ADCを再生成
gcloud auth application-default login

# ファイルをコピー
cp ~/.config/gcloud/application_default_credentials.json \
   backend/application_default_credentials.json

# Dockerを再起動
docker-compose restart backend
```

## 📊 動作確認

### バックエンドログの確認

```bash
# リアルタイムでログを確認
docker-compose logs -f backend
```

**正常な起動ログ:**
```
Server is running on http://localhost:8080
GraphQL Playground: http://localhost:8080/graphql
```

### 画像生成のテスト

```bash
# curlでテスト
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { generateImage(input: { text: \"test\" }) { imageUrl } }"
  }'
```

## 💰 コスト情報

### Imagen 2 の料金

| 操作 | 料金 |
|------|------|
| 1回の画像生成 (1024x1024) | 約$0.02 |
| 100回の生成 | 約$2.00 |

### コスト削減のヒント

1. **開発中は最小限に**: 必要な時だけ画像生成
2. **キャッシュを活用**: 同じプロンプトは再利用を検討
3. **使用量を監視**: GCPコンソールで定期的にチェック

```bash
# 使用量の確認
gcloud billing accounts list
```

## 🔍 デバッグ情報

### 実装の詳細

- **ファイル**: `backend/image/rest_client.go`
- **エンドポイント**: `https://us-central1-aiplatform.googleapis.com/v1/projects/{project}/locations/{location}/publishers/google/models/imagegeneration@002:predict`
- **認証**: OAuth2 Bearer token (ADC)
- **レスポンス**: Base64エンコードされたPNG画像

### APIリクエスト例

```json
{
  "instances": [
    {"prompt": "A dramatic fire scene"}
  ],
  "parameters": {
    "sampleCount": 1,
    "aspectRatio": "1:1",
    "negativePrompt": "blurry, low quality",
    "sampleImageSize": "1024"
  }
}
```

## 📚 参考リンク

- [Vertex AI - Imagen Documentation](https://cloud.google.com/vertex-ai/docs/generative-ai/image/overview)
- [Vertex AI REST API](https://cloud.google.com/vertex-ai/docs/reference/rest)
- [Authentication](https://cloud.google.com/docs/authentication/application-default-credentials)

## ✨ 次のステップ

画像生成が成功したら:

1. **異なるスタイルを試す**: MEME, REALISTIC, ILLUSTRATION, DRAMATIC
2. **プロンプトを調整**: より詳細な説明で品質向上
3. **Twitter投稿**: 生成した画像をSNSで共有
4. **統合テストを実行**: `RUN_INTEGRATION_TESTS=true go test ./image -v`

---

**問題が解決しない場合**: GitHubのIssuesで質問するか、バックエンドログの全文を共有してください。
