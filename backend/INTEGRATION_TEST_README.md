# Integration Tests

このディレクトリには、画像生成機能の統合テストが含まれています。

## 概要

統合テストは、実際のVertex AI APIを使用して、完全なフローをテストします。

- `image/integration_test.go` - Imagen API クライアントの統合テスト
- `graph/integration_test.go` - GraphQL API の統合テスト

## 実行方法

### 通常のテスト実行（統合テストはスキップ）

```bash
# すべてのテストを実行（統合テストはスキップされる）
make backend-check

# または
go test ./...
```

統合テストはデフォルトでスキップされます。これは以下の理由によります：
- 実際のGCP APIを呼び出すため、課金が発生する
- ネットワーク接続が必要
- 実行に時間がかかる（画像生成は5-15秒程度）

### 統合テストを実行

統合テストを実行するには、環境変数を設定します：

```bash
# 環境変数を設定
export RUN_INTEGRATION_TESTS=true
export GCP_PROJECT_ID=your-project-id
export GCP_LOCATION=us-central1  # オプション（デフォルト: us-central1）

# 画像生成の統合テストのみ実行
go test ./image -v -run Integration

# GraphQL APIの統合テストのみ実行
go test ./graph -v -run Integration

# すべての統合テストを実行
go test ./... -v -run Integration
```

## 前提条件

統合テストを実行するには、以下が必要です：

### 1. GCPプロジェクトのセットアップ

```bash
# Vertex AI APIを有効化
gcloud services enable aiplatform.googleapis.com

# 現在のプロジェクトIDを確認
gcloud config get-value project
```

### 2. 認証情報の設定

```bash
# Application Default Credentials (ADC) を設定
gcloud auth application-default login

# または、サービスアカウントキーを使用
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
```

### 3. 必要な権限

サービスアカウント（またはユーザーアカウント）に以下の権限が必要です：

- `roles/aiplatform.user` - Vertex AI の使用

## テスト内容

### image/integration_test.go

#### TestImageGenerationIntegration

1. **complete image generation flow**
   - クライアント作成
   - 画像生成
   - 結果の検証（PNG形式、データサイズなど）

2. **image generation with style options**
   - 異なるスタイル（realistic, illustration）での画像生成
   - オプション（style, aspectRatio）の動作確認

3. **multiple concurrent image generations**
   - 並行画像生成（3つ同時）
   - レート制限の確認

4. **error handling - invalid prompt**
   - 空のプロンプトでのエラーハンドリング

#### TestImageGenerationPerformance

1. **measure image generation latency**
   - ウォームアップリクエスト
   - パフォーマンス測定
   - メトリクスのログ出力

### graph/integration_test.go

#### TestGraphQLImageGenerationIntegration

1. **complete GraphQL generateImage flow**
   - Gemini + Imagen の完全な統合フロー
   - GraphQLリゾルバーの動作確認
   - Base64エンコードされた画像データの検証

2. **generateImage with different styles**
   - 各スタイル（MEME, REALISTIC など）のテスト

3. **generateImage with different aspect ratios**
   - 各アスペクト比（SQUARE など）のテスト

4. **error handling - empty text**
   - 入力検証のテスト

#### TestGraphQLImageGenerationWithTwitterIntegration

1. **prepare data for Twitter posting**
   - 画像生成からTwitter投稿準備までのフロー
   - データ構造の検証

## コスト管理

統合テストは実際のAPIを呼び出すため、コストが発生します。

### Imagen 3 の価格（2025年1月時点）

| 画像サイズ | 価格/画像 |
|-----------|----------|
| 512x512   | $0.020   |
| 1024x1024 | $0.040   |

### テスト実行時のコスト見積もり

統合テストを1回実行すると、約10-15枚の画像が生成されます：

- 開発環境（512x512）: $0.20 - $0.30 / 実行
- 本番環境（1024x1024）: $0.40 - $0.60 / 実行

### コスト削減のヒント

1. **必要な時だけ実行**: CI/CDでは統合テストをスキップ、手動実行のみ
2. **並行実行を制限**: レート制限とコストの両面から、並行数を3以下に
3. **小さい画像サイズ**: テストでは512x512を使用

## CI/CD での実行

### GitHub Actions の例

```yaml
name: Integration Tests

on:
  workflow_dispatch:  # 手動実行のみ
  schedule:
    - cron: '0 0 * * 0'  # 週1回（日曜日）のみ

jobs:
  integration-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21

      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v1
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - name: Run Integration Tests
        env:
          RUN_INTEGRATION_TESTS: true
          GCP_PROJECT_ID: ${{ secrets.GCP_PROJECT_ID }}
          GCP_LOCATION: us-central1
        run: |
          cd backend
          go test ./... -v -run Integration
```

## トラブルシューティング

### 1. 認証エラー

```
Error: google: could not find default credentials
```

**解決方法:**
```bash
gcloud auth application-default login
```

### 2. APIが無効

```
Error: Vertex AI API has not been used in project
```

**解決方法:**
```bash
gcloud services enable aiplatform.googleapis.com
```

### 3. 権限エラー

```
Error: Permission 'aiplatform.endpoints.predict' denied
```

**解決方法:**
```bash
# IAMロールを確認・追加
gcloud projects add-iam-policy-binding PROJECT_ID \
  --member="user:YOUR_EMAIL" \
  --role="roles/aiplatform.user"
```

### 4. レート制限エラー

```
Error: Quota exceeded
```

**解決方法:**
- 並行実行数を減らす
- リトライ間隔を空ける
- GCPコンソールでクォータを確認

## ベストプラクティス

1. **ローカル開発では慎重に実行**
   - 統合テストは必要な時だけ実行
   - 単体テストで十分カバーできる箇所は統合テスト不要

2. **テストデータは最小限に**
   - シンプルなプロンプトを使用
   - 必要最小限の画像生成に留める

3. **タイムアウトの設定**
   - 画像生成は時間がかかるため、適切なタイムアウトを設定

4. **エラーログの保存**
   - 失敗時のデバッグに役立つログを残す

## 参考リンク

- [Vertex AI - Imagen Documentation](https://cloud.google.com/vertex-ai/docs/generative-ai/image/overview)
- [Google Cloud Authentication](https://cloud.google.com/docs/authentication)
- [Vertex AI Pricing](https://cloud.google.com/vertex-ai/pricing)
