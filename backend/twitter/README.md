# Twitter API クライアント

炎上シミュレーターのTwitter投稿機能を提供するGoパッケージ。

## 機能

- ✅ テキスト投稿（OAuth 1.0a認証）
- ✅ 画像付き投稿（Media Upload API実装済み）
- ✅ ハッシュタグと免責文言の自動追加
- ✅ 280文字制限のバリデーション
- ✅ モック/本番モード自動切り替え

## 使用方法

### クライアントの初期化

```go
import "github.com/Tattsum/enjo/backend/twitter"

client, err := twitter.NewClient(
    apiKey,
    apiSecret,
    accessToken,
    accessTokenSecret,
)
if err != nil {
    log.Fatal(err)
}
```

### テキスト投稿

```go
ctx := context.Background()
result, err := client.PostTweet(ctx, "投稿テキスト")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Posted: %s\n", result.URL)
```

### 画像付き投稿

```go
ctx := context.Background()
imageData := []byte{...} // PNG/JPEG画像データ

result, err := client.PostTweetWithImage(
    ctx,
    "投稿テキスト",
    imageData,
    twitter.WithHashtag(),
    twitter.WithDisclaimer(),
)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Posted with image: %s\n", result.URL)
```

## オプション

### WithHashtag()

ツイートに `#炎上シミュレーター` ハッシュタグを追加します。

### WithDisclaimer()

ツイートの末尾に `※炎上シミュレーターで生成` 免責文言を追加します。

## テスト

### ユニットテスト

```bash
go test ./twitter -v
```

### 統合テスト（実際のTwitter APIを呼び出す）

```bash
# 環境変数を設定
export RUN_TWITTER_INTEGRATION_TESTS=true
export TWITTER_API_KEY=your_api_key
export TWITTER_API_SECRET=your_api_secret
export TWITTER_ACCESS_TOKEN=your_access_token
export TWITTER_ACCESS_TOKEN_SECRET=your_access_token_secret

# 統合テストを実行
go test ./twitter -v -run Integration
```

**注意**: 統合テストは実際のTwitterに投稿されます。テストアカウントでの実行を推奨します。

## アーキテクチャ

### Media Upload API実装

go-twitterライブラリはMedia Upload APIをサポートしていないため、カスタム実装を提供しています。

**実装方式**:
- **エンドポイント**: `https://upload.twitter.com/1.1/media/upload.json`
- **認証**: OAuth 1.0a（既存のhttpClientを再利用）
- **リクエスト**: `application/x-www-form-urlencoded`
- **パラメータ**:
  - `media_data`: Base64エンコード画像
  - `media_category`: "tweet_image"
- **レスポンス**: `media_id_string`を取得

### モックモード

テスト環境では自動的にモックモードに切り替わります:

- APIキーが `test-api-key` または `test-key` の場合
- 実際のHTTPリクエストを送信せず、モックレスポンスを返す
- テストが高速に実行される

## ファイル構成

```
twitter/
├── client.go           # メインクライアント実装
├── client_test.go      # ユニットテスト（PostTweet）
├── media_test.go       # ユニットテスト（Media Upload）
├── integration_test.go # 統合テスト（実API呼び出し）
└── README.md          # このファイル
```

## 依存関係

- `github.com/dghubble/go-twitter/twitter` - Twitter APIクライアント（基本機能）
- `github.com/dghubble/oauth1` - OAuth 1.0a認証

## エラーハンドリング

### 一般的なエラー

- `"tweet text cannot be empty"` - 空のテキスト
- `"tweet text exceeds 280 characters"` - 文字数超過
- `"image data cannot be empty"` - 空の画像データ
- `"all Twitter API credentials are required"` - 認証情報不足

### Media Upload エラー

- `"failed to upload media"` - メディアアップロード失敗
- `"media upload failed with status XXX"` - HTTPエラー

### Twitter API エラー

- 401 Unauthorized - 認証失敗
- 403 Forbidden - 重複ツイートなど
- 429 Too Many Requests - レート制限超過

## パフォーマンス

### ベンチマーク

```bash
go test -bench=. ./twitter
```

### カバレッジ

```bash
go test -coverprofile=coverage.out ./twitter
go tool cover -html=coverage.out
```

現在のカバレッジ: **54.3%**

## 制限事項

### Twitter API制限

- ツイート投稿: 300ツイート/3時間（ユーザーごと）
- メディアアップロード: 制限あり（詳細はTwitter APIドキュメント参照）

### 実装上の制限

- 動画アップロードは未対応（画像のみ）
- 複数画像の同時アップロード未対応（1枚のみ）
- リツイート、引用ツイート機能未対応

## 今後の拡張

- [ ] レート制限対策（Exponential Backoff）
- [ ] リトライ機能
- [ ] 動画アップロード対応
- [ ] 複数画像対応
- [ ] より詳細なエラーメッセージ

## ライセンス

このプロジェクトのライセンスに従います。

## 関連ドキュメント

- [FEATURE_TWITTER_POST.md](../../docs/FEATURE_TWITTER_POST.md) - 詳細な機能仕様
- [Twitter API Documentation](https://developer.x.com/en/docs/twitter-api)
- [OAuth 1.0a Documentation](https://developer.x.com/en/docs/authentication/oauth-1-0a)
