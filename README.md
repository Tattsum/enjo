# 🔥 炎上シミュレーター (Enjo Simulator)

SNS投稿を「炎上しやすい表現」に変換し、予想されるリプライを生成するシミュレーター（教育・エンターテインメント目的）

## ⚠️ 注意事項

このツールは以下の目的のために作成されています:

- **教育目的**: 炎上のメカニズムを理解する
- **エンターテインメント**: 表現の違いによる印象の変化を体験する
- **コミュニケーション学習**: 誤解を招きやすい表現を学ぶ

**実際のSNSでの使用や、他者への攻撃・嫌がらせを目的とした使用は厳禁です。**

## 🎯 プロジェクト概要

通常の投稿文を入力すると、Google Gemini APIを使用して「炎上度レベル」に応じた表現に変換します。さらに、予想されるリプライ（正論批判、揚げ足取り、的外れな批判、過剰擁護）を生成します。

### 主な機能

1. **テキスト変換**: 炎上度レベル(1-5)に応じた表現変換
2. **リプライ生成**: 4種類の典型的なリプライパターンを自動生成
3. **比較表示**: 元の投稿と変換後を並べて表示
4. **説明生成**: なぜ炎上しやすいのかの解説
5. **Twitter投稿**: 生成したテキストをTwitter/𝕏に直接投稿（オプション）

## 🛠️ 技術スタック

### バックエンド
- **Go 1.25**
- **gqlgen** - GraphQL サーバー
- **Google Vertex AI (Gemini)** - AI テキスト生成
- **Twitter API v2** - SNS投稿機能（オプション）
- **chi** - HTTP ルーター
- **Air** - ホットリロード

### フロントエンド
- **Next.js 15**
- **React 19**
- **TypeScript 5.6**
- **Apollo Client** - GraphQL クライアント
- **Tailwind CSS** - スタイリング

### 開発環境
- **Docker & Docker Compose**
- **golangci-lint** - Go Linter
- **Make** - タスク自動化

## 📁 ディレクトリ構成

```text
enjo/
├── docker-compose.yml
├── Makefile
├── PROJECT_RULES.md        # 開発ルール (TDD)
├── README.md
├── backend/
│   ├── Dockerfile
│   ├── .air.toml           # ホットリロード設定
│   ├── .golangci.yml       # Lint設定
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── graph/
│   │   ├── schema.graphqls
│   │   ├── resolver.go
│   │   └── generated/
│   ├── gemini/
│   │   └── client.go
│   └── twitter/
│       └── client.go
└── frontend/
    ├── Dockerfile
    ├── .env.local.example
    ├── package.json
    ├── tsconfig.json
    ├── next.config.js
    └── src/
        ├── app/
        │   ├── layout.tsx
        │   └── page.tsx
        ├── components/
        │   ├── TextInput.tsx
        │   ├── LevelSlider.tsx
        │   ├── ResultDisplay.tsx
        │   ├── ReplyList.tsx
        │   └── TwitterPostButton.tsx
        └── lib/
            └── graphql/
                ├── client.ts
                └── queries.ts
```

## 🚀 セットアップ手順

### 前提条件

- Docker & Docker Compose
- Google Cloud Platform アカウント（Vertex AI用）
- Twitter Developer アカウント（オプション、Twitter投稿機能を使う場合）

### 1. リポジトリのクローン

```bash
git clone https://github.com/Tattsum/enjo.git
cd enjo
```

### 2. 環境変数の設定

```bash
# Makefileを使った自動セットアップ
make setup

# または手動で設定
cp backend/.env.example backend/.env
cp frontend/.env.local.example frontend/.env.local
```

### 3. Google Cloud Platform (Vertex AI) の設定

#### GCPプロジェクトとVertex AIの有効化

```bash
# Google Cloud SDK にログイン
gcloud auth login
gcloud config set project YOUR_PROJECT_ID

# Vertex AI API を有効化
gcloud services enable aiplatform.googleapis.com

# Application Default Credentials (ADC) を設定
gcloud auth application-default login
```

詳細な手順は [docs/SETUP_API_KEY.md](docs/SETUP_API_KEY.md) を参照してください。

#### 環境変数への設定

`backend/.env` を編集:

```env
GCP_PROJECT_ID=your_gcp_project_id_here
GCP_LOCATION=us-central1
PORT=8080
```

### 4. Twitter API の設定（オプション）

Twitter投稿機能を使用する場合のみ設定してください。

#### Twitter Developer Portalでアプリを作成

1. [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard) にアクセス
2. 「Create App」をクリックしてアプリを作成
3. App Permissions を「Read and Write」に設定
4. API Key & Secret、Access Token & Secret を取得

#### 環境変数への設定

`backend/.env` に追記:

```env
# Twitter API Configuration (Optional)
TWITTER_API_KEY=your_twitter_api_key_here
TWITTER_API_SECRET=your_twitter_api_secret_here
TWITTER_ACCESS_TOKEN=your_twitter_access_token_here
TWITTER_ACCESS_TOKEN_SECRET=your_twitter_access_token_secret_here
```

**注意**: Twitter API認証情報が設定されていない場合、Twitter投稿機能は無効化されますが、その他の機能は正常に動作します。

詳細は [docs/FEATURE_TWITTER_POST.md](docs/FEATURE_TWITTER_POST.md) を参照してください。

### 5. Docker Composeで起動

```bash
# すべてのサービスを起動
make docker-up

# または
docker-compose up --build
```

### 6. アクセス

- **フロントエンド**: <http://localhost:3000>
- **GraphQL Playground**: <http://localhost:8080/graphql>
- **バックエンドヘルスチェック**: <http://localhost:8080/health>

## 🧪 開発ワークフロー (TDD)

このプロジェクトは **t-wada スタイル**のテスト駆動開発(TDD)に従います。

### Red-Green-Refactor サイクル

```text
1. Red:      失敗するテストを書く
2. Green:    テストが通る最小限のコードを書く
3. Refactor: コードをきれいにする
4. Repeat:   次の機能へ
```

### Makeコマンド

```bash
# ヘルプを表示
make help

# すべてのチェック (fmt → lint → test)
make check

# バックエンドのみ
make backend-check
make backend-test
make backend-lint
make backend-fmt

# フロントエンドのみ
make frontend-check
make frontend-test
make frontend-lint
make frontend-fmt
```

詳細は [PROJECT_RULES.md](PROJECT_RULES.md) を参照してください。

## 📝 使い方

### 基本的な使い方

1. テキストエリアに普通の投稿を入力
2. 炎上度スライダー(1-5)を調整
3. 「🔥 炎上化する」ボタンをクリック
4. 元の投稿と変換後を比較
5. 「💬 リプライを生成」で予想されるリプライを表示
6. （オプション）「𝕏 Xに投稿」ボタンでTwitterに投稿

### 炎上度レベル

- **レベル1**: 少し配慮に欠ける表現
- **レベル2**: 誤解を招きやすい表現
- **レベル3**: 明確に批判されそうな表現
- **レベル4**: かなり問題がある表現
- **レベル5**: 炎上確実な表現

### リプライタイプ

- **正論で批判**: 論理的な批判
- **揚げ足取り**: 些細な点を指摘
- **的外れな批判**: 本質とは関係ない批判
- **過剰擁護**: 過度に擁護する意見

## 🔧 ローカル開発（Docker なし）

### バックエンド

```bash
cd backend

# 依存関係のインストール
go mod download

# 開発サーバー起動 (ホットリロード)
air

# または通常起動
go run main.go
```

### フロントエンド

```bash
cd frontend

# 依存関係のインストール
npm install

# 開発サーバー起動
npm run dev
```

## 🧪 テスト

```bash
# すべてのテスト
make test

# バックエンドのみ
make backend-test

# フロントエンドのみ
make frontend-test
```

## 📊 GraphQL API

### ヘルスチェック

```graphql
query {
  health
}
```

### テキスト変換

```graphql
mutation {
  generateInflammatoryText(input: {
    originalText: "今日はいい天気ですね"
    level: 3
  }) {
    inflammatoryText
    explanation
  }
}
```

### リプライ生成

```graphql
mutation {
  generateReplies(text: "変換されたテキスト") {
    id
    type
    content
  }
}
```

### Twitter投稿（オプション）

```graphql
mutation {
  postToTwitter(input: {
    text: "投稿するテキスト"
    addHashtag: true
    addDisclaimer: true
  }) {
    success
    tweetId
    tweetUrl
    errorMessage
  }
}
```

## 🤝 コントリビューション

プルリクエストを歓迎します！

1. このリポジトリをフォーク
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. TDDに従って開発 (`make check` でテスト)
4. コミット (`git commit -m 'Add amazing feature'`)
5. プッシュ (`git push origin feature/amazing-feature`)
6. プルリクエストを作成

### コーディング規約

- [PROJECT_RULES.md](PROJECT_RULES.md) に従ってください
- すべてのコードにテストを書く
- `make check` が通ることを確認

## 📄 ライセンス

MIT License

Copyright (c) 2025 Tattsum

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## 🙏 謝辞

- [Google Vertex AI (Gemini)](https://cloud.google.com/vertex-ai) - AI テキスト生成
- [Twitter API v2](https://developer.twitter.com/en/docs/twitter-api) - SNS投稿機能
- [gqlgen](https://gqlgen.com/) - GraphQL サーバー
- [Next.js](https://nextjs.org/) - React フレームワーク

## 📮 お問い合わせ

質問や提案がある場合は、[Issues](https://github.com/Tattsum/enjo/issues) を作成してください。

---

**Remember**: このツールは教育・エンターテインメント目的です。実際のSNSで他者を傷つける目的での使用は絶対にしないでください。
