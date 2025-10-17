# 画像自動生成機能 (Image Generation Feature)

## 概要

炎上投稿に合わせた画像を自動生成し、テキストと一緒にTwitter/𝕏に投稿できる機能。

## 目的

- 炎上シミュレーションをより視覚的に表現
- SNS投稿のエンゲージメントを高める
- 自動化されたコンテンツ生成のデモンストレーション

## 技術選定

### 画像生成API: Google Vertex AI - Imagen 3

**選定理由:**
- 既にVertex AI (Gemini)を使用しているため、認証・インフラが統一できる
- 同一のGCPプロジェクトで管理可能
- 高品質な画像生成が可能
- 日本語プロンプトのサポート
- コスト管理が容易

**代替案との比較:**

| 項目 | Imagen 3 (推奨) | DALL-E 3 | Stable Diffusion |
|------|----------------|----------|-----------------|
| 認証統合 | ✅ Vertex AI統合 | ❌ 別途APIキー | ⚠️ 自前ホスティング |
| 日本語対応 | ✅ ネイティブ | ⚠️ 英訳が必要 | ⚠️ 英訳が必要 |
| コスト | 中程度 | 高い | 低い（インフラ費用別） |
| 品質 | 高い | 非常に高い | 中〜高い |
| レイテンシ | 5-15秒 | 10-30秒 | 変動大 |

## アーキテクチャ

### システムフロー

```
1. ユーザーが炎上テキストを生成
   ↓
2. Gemini APIで画像プロンプトを生成
   ↓
3. Imagen APIで画像を生成
   ↓
4. 生成した画像をフロントエンドに表示
   ↓
5. (オプション) Twitter APIで画像とテキストを投稿
```

### コンポーネント構成

```
backend/
├── image/
│   ├── client.go           # Imagen APIクライアント
│   ├── client_test.go      # テスト
│   └── prompt.go           # 画像プロンプト生成ロジック
├── graph/
│   └── schema.graphqls     # GraphQLスキーマ拡張
└── twitter/
    └── client.go           # 画像付き投稿対応

frontend/
└── src/
    └── components/
        ├── ImageGenerator.tsx        # 画像生成UI
        ├── ImagePreview.tsx          # 画像プレビュー
        └── TwitterPostButton.tsx     # 画像付き投稿対応
```

## API設計

### GraphQL Schema拡張

```graphql
# 画像生成のミューテーション
mutation {
  generateImage(input: GenerateImageInput!): GenerateImageResult!
}

input GenerateImageInput {
  text: String!              # 炎上テキスト
  style: ImageStyle          # 画像スタイル (オプション)
  aspectRatio: AspectRatio   # アスペクト比 (オプション)
}

enum ImageStyle {
  REALISTIC      # リアル調
  ILLUSTRATION   # イラスト調
  MEME          # ミーム風
  DRAMATIC      # ドラマチック
}

enum AspectRatio {
  SQUARE        # 1:1 (Twitter最適)
  LANDSCAPE     # 16:9
  PORTRAIT      # 9:16
}

type GenerateImageResult {
  imageUrl: String!          # 生成された画像のURL
  prompt: String!            # 使用したプロンプト
  generatedAt: String!       # 生成日時
}

# Twitter投稿のミューテーション拡張
mutation {
  postToTwitter(input: TwitterPostInput!): TwitterPostResult!
}

input TwitterPostInput {
  text: String!
  imageUrl: String           # 画像URL (オプション)
  addHashtag: Boolean
  addDisclaimer: Boolean
}
```

### バックエンドAPI

#### 1. Imagen クライアント

```go
// backend/image/client.go

package image

import (
    "context"
    "cloud.google.com/go/vertexai/genai"
)

type Client struct {
    client    *genai.Client
    projectID string
    location  string
}

// NewClient creates a new Imagen API client
func NewClient(ctx context.Context, projectID, location string) (*Client, error)

// GenerateImage generates an image based on the prompt
func (c *Client) GenerateImage(
    ctx context.Context,
    prompt string,
    options ...ImageOption
) (*ImageResult, error)

type ImageResult struct {
    ImageData   []byte    // 画像データ (PNG)
    ImageURL    string    // GCS URL (保存した場合)
    Prompt      string    // 使用したプロンプト
    GeneratedAt time.Time
}

type ImageOption func(*imageOptions)

func WithStyle(style string) ImageOption
func WithAspectRatio(ratio string) ImageOption
func WithSize(width, height int) ImageOption
```

#### 2. プロンプト生成ロジック

```go
// backend/image/prompt.go

// GenerateImagePrompt generates an image generation prompt from inflammatory text
func GenerateImagePrompt(ctx context.Context, geminiClient *gemini.Client, text string) (string, error) {
    // Geminiを使って炎上テキストから画像プロンプトを生成
    prompt := fmt.Sprintf(`
以下の炎上投稿に合わせた、視覚的にインパクトのある画像のプロンプトを生成してください。

【投稿】
%s

【要件】
- 投稿の雰囲気を視覚的に表現
- 炎のモチーフを含める
- SNS映えする構図
- ミーム的な要素
- 日本のネット文化に馴染む表現

画像生成プロンプト（英語）のみを出力してください。
`, text)

    return geminiClient.GenerateContent(ctx, prompt)
}
```

#### 3. Twitter クライアント拡張

```go
// backend/twitter/client.go

// PostTweetWithImage posts a tweet with an image
func (c *Client) PostTweetWithImage(
    ctx context.Context,
    text string,
    imageData []byte,
    options ...TweetOption
) (*TweetResult, error) {
    // 1. メディアをアップロード
    mediaID, err := c.uploadMedia(ctx, imageData)
    if err != nil {
        return nil, err
    }

    // 2. ツイートを投稿
    return c.postTweetWithMediaID(ctx, text, mediaID, options...)
}
```

## フロントエンド設計

### コンポーネント

#### 1. ImageGenerator コンポーネント

```typescript
// frontend/src/components/ImageGenerator.tsx

interface ImageGeneratorProps {
  inflammatoryText: string;
  onImageGenerated?: (imageUrl: string) => void;
}

export function ImageGenerator({ inflammatoryText, onImageGenerated }: ImageGeneratorProps) {
  const [generateImage] = useMutation(GENERATE_IMAGE);
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  // 画像生成処理
  // プレビュー表示
  // スタイル選択UI
}
```

#### 2. ImagePreview コンポーネント

```typescript
// frontend/src/components/ImagePreview.tsx

interface ImagePreviewProps {
  imageUrl: string;
  prompt?: string;
  onDownload?: () => void;
  onRegenerate?: () => void;
}

export function ImagePreview({ imageUrl, prompt, onDownload, onRegenerate }: ImagePreviewProps) {
  // 画像表示
  // ダウンロードボタン
  // 再生成ボタン
}
```

### GraphQLクエリ

```typescript
// frontend/src/lib/graphql/queries.ts

export const GENERATE_IMAGE = gql`
  mutation GenerateImage($input: GenerateImageInput!) {
    generateImage(input: $input) {
      imageUrl
      prompt
      generatedAt
    }
  }
`;

export const POST_TO_TWITTER_WITH_IMAGE = gql`
  mutation PostToTwitterWithImage($input: TwitterPostInput!) {
    postToTwitter(input: $input) {
      success
      tweetId
      tweetUrl
      errorMessage
    }
  }
`;
```

## ユーザーフロー

### 基本フロー

1. ユーザーが炎上テキストを生成
2. 「🎨 画像を生成」ボタンをクリック
3. (オプション) 画像スタイルを選択
4. 画像が生成され、プレビュー表示
5. 必要に応じて再生成可能
6. 「𝕏 画像付きで投稿」ボタンでTwitterに投稿

### UI配置イメージ

```
���─────────────────────────────────────┐
│ 炎上シミュレーター                    │
├─────────────────────────────────────┤
│ [テキスト入力エリア]                  │
│ [炎上度スライダー: 1-5]               │
│ [🔥 炎上化する] ボタン                 │
├─────────────────────────────────────┤
│ 【結果表示】                          │
│ ┌─────────┐  ┌─────────┐           │
│ │ 元の投稿 │  │変換後   │           │
│ └─────────┘  └─────────┘           │
├─────────────────────────────────────┤
│ [🎨 画像を生成] ボタン ← NEW!         │
│                                     │
│ ┌─────────────────────────┐         │
│ │  生成された画像           │ ← NEW!  │
│ │  [プレビュー表示]        │         │
│ └─────────────────────────┘         │
│ [↻ 再生成] [⬇ ダウンロード]          │
├─────────────────────────────────────┤
│ [💬 リプライを生成]                   │
│ [𝕏 Xに投稿] / [𝕏 画像付きで投稿] ← NEW!│
└─────────────────────────────────────┘
```

## GCP設定

### 必要なAPI

```bash
# Vertex AI - Imagen API を有効化
gcloud services enable aiplatform.googleapis.com

# 既に有効化されているはずだが、念のため確認
gcloud services list --enabled | grep aiplatform
```

### IAM権限

プロジェクトのサービスアカウントに以下の権限が必要:

- `roles/aiplatform.user` - Vertex AI の使用
- `roles/storage.objectCreator` - GCS への画像保存（オプション）

### 環境変数

`backend/.env` に追加:

```env
# Image Generation Configuration
GCP_PROJECT_ID=your_gcp_project_id_here
GCP_LOCATION=us-central1
ENABLE_IMAGE_GENERATION=true

# Optional: GCS bucket for storing generated images
GCS_BUCKET_NAME=enjo-generated-images
```

## コスト見積もり

### Imagen 3 価格 (2025年1月時点)

| 項目 | 価格 |
|------|------|
| 画像生成 (512x512) | $0.020 / 画像 |
| 画像生成 (1024x1024) | $0.040 / 画像 |

### 月間コスト試算

| 使用量 | コスト (512x512) | コスト (1024x1024) |
|--------|------------------|-------------------|
| 100画像/月 | $2.00 | $4.00 |
| 1,000画像/月 | $20.00 | $40.00 |
| 10,000画像/月 | $200.00 | $400.00 |

**推奨:** 開発環境では512x512、本番環境では1024x1024を使用

## セキュリティ考慮事項

### 画像生成の制限

1. **レート制限**
   - ユーザーあたり: 10画像/時間
   - IP あたり: 50画像/時間

2. **コンテンツフィルタリング**
   - Imagen のセーフティフィルターを有効化
   - 不適切なプロンプトの検出

3. **画像保存**
   - 生成画像は24時間後に自動削除
   - GCS のライフサイクル管理を使用

### プライバシー

- 生成した画像にユーザー識別情報を含めない
- 画像URLは予測不可能なランダム文字列を使用

## 実装ステップ (TDD)

### Phase 1: バックエンド基盤 ✅ 完了

1. **Imagen クライアントの実装** ✅
   - [x] テスト作成: `backend/image/client_test.go`
   - [x] 実装: `backend/image/client.go`
   - [x] 認証・初期化テスト
   - [x] 画像生成テスト（統合テストは分離済み）
   - [x] すべてのテストがパス（`make backend-check`）

2. **プロンプト生成ロジック** ✅
   - [x] テスト作成: `backend/image/prompt_test.go`
   - [x] 実装: `backend/image/prompt.go`
   - [x] Gemini連携テスト（モック使用）

### Phase 2: GraphQL API ✅ 完了

3. **GraphQL スキーマ拡張** ✅
   - [x] スキーマ定義: `backend/graph/schema.graphqls`
     - `GenerateImageInput`, `GenerateImageResult` 型定義
     - `ImageStyle`, `AspectRatio` Enum定義
     - `generateImage` ミューテーション追加
   - [x] gqlgen コード生成実行
   - [x] リゾルバーテスト: `backend/graph/resolver_test.go`
     - `TestMutationResolver_GenerateImage` 追加（4テストケース）
     - モック実装（MockGeminiClient, MockImageClient）
   - [x] リゾルバー実装: `backend/graph/schema.resolvers.go`
     - `GenerateImage` リゾルバー実装
     - ヘルパー関数追加（`generateImagePromptFromText`, `createImageDataURL`, `getCurrentTimestamp`）
   - [x] インターフェース定義: `backend/graph/resolver.go`
     - `GeminiClient` に `GenerateContent` メソッド追加
     - `ImageClient` インターフェース追加
   - [x] 統合
     - `backend/image/adapter.go` でアダプターパターン実装
     - `backend/gemini/client.go` に `GenerateContent` 公開メソッド追加
     - `backend/main.go` でimageClient初期化と注入
     - すべてのテストがパス（`make backend-check`）
   - [x] コードカバレッジ: 63.6%

### Phase 3: Twitter連携 ✅ 完了

4. **Twitter画像投稿** ✅
   - [x] テスト作成: `backend/twitter/media_test.go`
   - [x] メディアアップロード実装
   - [x] 画像付き投稿実装
   - [x] GraphQLスキーマ拡張（TwitterPostInputにimageUrlフィールド追加）
   - [x] PostToTwitterリゾルバーの画像対応
   - [x] すべてのテストがパス（`make backend-check`）
   - [x] カバレッジ: twitter 89.8%

### Phase 4: フロントエンド

5. **ImageGenerator コンポーネント**
   - [ ] テスト作成: `frontend/src/components/ImageGenerator.test.tsx`
   - [ ] コンポーネント実装
   - [ ] GraphQLクエリ統合

6. **ImagePreview コンポーネント**
   - [ ] テスト作成: `frontend/src/components/ImagePreview.test.tsx`
   - [ ] コンポーネント実装

7. **統合**
   - [ ] ResultDisplay コンポーネントに統合
   - [ ] TwitterPostButton の画像対応

### Phase 5: E2Eテスト

8. **統合テスト**
   - [ ] バックエンド統合テスト
   - [ ] フロントエンド E2E テスト
   - [ ] パフォーマンステスト

## パフォーマンス最適化

### 画像生成の最適化

1. **非同期処理**
   - 画像生成は時間がかかるため、非同期で処理
   - WebSocketまたはポーリングで進捗通知

2. **キャッシング**
   - 同じテキストの画像は再利用
   - Redis/Memcached でキャッシュ

3. **CDN配信**
   - Cloud CDN で画像を配信
   - レイテンシ削減

## テスト計画

### 単体テスト

- [ ] Imagen クライアントのテスト
- [ ] プロンプト生成ロジックのテスト
- [ ] Twitter メディアアップロードのテスト
- [ ] React コンポーネントのテスト

### 統合テスト

- [ ] GraphQL API のエンドツーエンドテスト
- [ ] 画像生成フローのテスト
- [ ] Twitter投稿フローのテスト

### E2Eテスト

- [ ] ユーザーフロー全体のテスト
- [ ] エラーハンドリングのテスト

## リリース計画

### v1.0 (MVP)

- [x] テキスト変換機能
- [x] リプライ生成機能
- [x] Twitter投稿機能 (テキストのみ)

### v1.1 (画像生成機能)

- [ ] 画像自動生成機能
- [ ] 画像プレビュー
- [ ] 画像付きTwitter投稿

### v1.2 (将来の拡張)

- [ ] 複数スタイルの画像生成
- [ ] 画像編集機能
- [ ] ギャラリー機能

## トラブルシューティング

### よくある問題

#### 1. Imagen API が有効化されていない

```bash
gcloud services enable aiplatform.googleapis.com
```

#### 2. 認証エラー

```bash
# Application Default Credentials を再設定
gcloud auth application-default login
```

#### 3. 画像生成が遅い

- 最初の呼び出しは遅い (コールドスタート)
- 512x512 サイズを使用して高速化
- バッチ処理を検討

## 実装の詳細

### Phase 1 & 2 実装サマリー（2025-10-17完了）

#### 実装したファイル

**Phase 1: バックエンド基盤**
- `backend/image/client.go` - Imagen APIクライアント
- `backend/image/client_test.go` - クライアントの単体テスト
- `backend/image/prompt.go` - プロンプト生成ロジック
- `backend/image/prompt_test.go` - プロンプト生成のテスト

**Phase 2: GraphQL API統合**
- `backend/graph/schema.graphqls` - GraphQLスキーマ拡張
- `backend/graph/resolver.go` - ImageClientインターフェース定義
- `backend/graph/schema.resolvers.go` - GenerateImageリゾルバー実装
- `backend/graph/helpers.go` - ヘルパー関数
- `backend/graph/resolver_test.go` - リゾルバーのテスト
- `backend/image/adapter.go` - ImageClientアダプター
- `backend/gemini/client.go` - GenerateContentメソッド追加
- `backend/main.go` - imageClient統合
- `backend/main_test.go` - メインのテスト更新

#### TDD準拠

すべてのコードはTDD（Red-Green-Refactor）サイクルに従って実装：
1. **Red**: テストを先に書いて失敗を確認
2. **Green**: 最小限の実装でテストをパス
3. **Refactor**: コードをリファクタリング
4. **Check**: `make backend-check`で品質確認

#### テスト結果

```bash
$ make backend-check
✅ golangci-lint: 0 issues
✅ すべてのテスト: PASS
✅ カバレッジ:
   - graph: 63.6%
   - image: 53.1%
   - backend: 20.8%
```

#### 注意事項

- 統合テストは`t.Skip()`でマークし、通常のテスト実行ではスキップ
- Vertex AI APIの非推奨警告あり（2026年6月24日まで使用可能）
- 将来的にGoogle GenAI Go SDKへの移行を検討

## 参考リンク

- [Vertex AI - Imagen Documentation](https://cloud.google.com/vertex-ai/docs/generative-ai/image/overview)
- [Twitter API v2 - Media Upload](https://developer.twitter.com/en/docs/twitter-api/v1/media/upload-media/overview)
- [GraphQL Best Practices](https://graphql.org/learn/best-practices/)
- [Google GenAI Go SDK](https://pkg.go.dev/google.golang.org/genai) - 将来の移行先

## まとめ

この機能により、炎上シミュレーターは以下の点で強化されます:

1. **視覚的なインパクト**: テキストだけでなく画像も自動生成
2. **SNS最適化**: 画像付き投稿でエンゲージメント向上
3. **自動化**: Gemini → Imagen → Twitter の完全自動フロー
4. **拡張性**: 将来的な画像編集・スタイル選択への拡張が容易

### 現在の実装状況（2025-10-17更新）

- ✅ **Phase 1**: バックエンド基盤（完了）
  - Imagenクライアント実装
  - プロンプト生成ロジック実装
  - すべてのテストがパス

- ✅ **Phase 2**: GraphQL API統合（完了）
  - GraphQLスキーマ拡張（generateImageミューテーション）
  - リゾルバー実装
  - main.goへの統合
  - すべてのテストがパス（`make backend-check`）

- ✅ **Phase 3**: Twitter連携（完了）
  - `backend/twitter/media_test.go` - メディアアップロードテスト作成
  - `backend/twitter/client.go` - 画像付き投稿機能実装
    - `uploadMedia` メソッド
    - `postTweetWithMediaID` メソッド
    - `PostTweetWithImage` 公開メソッド
  - GraphQLスキーマ拡張（TwitterPostInputにimageUrlフィールド追加）
  - PostToTwitterリゾルバーの画像対応
  - `backend/graph/helpers.go` - `extractImageDataFromURL` 関数追加
  - すべてのテストがパス（`make backend-check`）
  - カバレッジ: twitter 89.8%, graph 47.3%

- ⏳ **Phase 4**: フロントエンド（未着手）
  - ImageGeneratorコンポーネント
  - ImagePreviewコンポーネント
  - GraphQLクエリ統合

TDDに従い、小さく作って育てる方針で段階的に実装しました。Phase 1, 2, 3は完全にテスト駆動（Red-Green-Refactor）で実装され、すべてのテストがパスしています。
