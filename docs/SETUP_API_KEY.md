# API キー取得とセットアップ方法

このドキュメントでは、Enjo Simulatorで使用するAPI（Gemini、Imagen、Twitter）のセットアップ方法を説明します。

## 目次

1. [Gemini API (Generative Language API)](#gemini-api-generative-language-api)
2. [Imagen API (Vertex AI)](#imagen-api-vertex-ai)
3. [Twitter API](#twitter-api)

---

## Gemini API (Generative Language API)

### 前提条件

- Google Cloud SDK (gcloud) がインストールされていること
- Google Cloud プロジェクトが作成されていること
- Generative Language API が有効化されていること

## 手順

### 1. gcloud のインストール確認

```bash
gcloud --version
```

インストールされていない場合は、[Google Cloud SDK のインストール](https://cloud.google.com/sdk/docs/install)を参照してください。

### 2. gcloud の初期化とログイン

```bash
# gcloud の初期化
gcloud init

# または既にプロジェクトがある場合
gcloud auth login
gcloud config set project YOUR_PROJECT_ID
```

### 3. Generative Language API の有効化

```bash
# API を有効化
gcloud services enable generativelanguage.googleapis.com
```

### 4. API キーの作成

Google Cloud Console または gcloud コマンドで API キーを作成します。

#### 方法1: gcloud コマンドで作成

```bash
# API キーを作成
gcloud alpha services api-keys create \
  --display-name="Enjo Simulator API Key" \
  --api-target=service=generativelanguage.googleapis.com

# 作成されたキーの一覧を表示
gcloud alpha services api-keys list

# 特定のキーの詳細を取得 (KEY_ID は上記コマンドで確認)
gcloud alpha services api-keys get-key-string KEY_ID
```

#### 方法2: Google Cloud Console で作成

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. プロジェクトを選択
3. 「APIとサービス」→「認証情報」を選択
4. 「認証情報を作成」→「APIキー」を選択
5. 作成されたAPIキーをコピー

### 5. API キーの制限設定（推奨）

セキュリティのため、APIキーに制限を設定することを推奨します。

```bash
# APIキーを特定のAPIに制限
gcloud alpha services api-keys update KEY_ID \
  --api-target=service=generativelanguage.googleapis.com

# IPアドレス制限を追加する場合
gcloud alpha services api-keys update KEY_ID \
  --allowed-ips=YOUR_IP_ADDRESS
```

### 6. .env ファイルの作成

取得したAPIキーをプロジェクトの `.env` ファイルに設定します。

```bash
# backend ディレクトリに移動
cd backend

# .env.example をコピー
cp .env.example .env

# エディタで .env を編集
# GEMINI_API_KEY に取得したキーを設定
```

`.env` ファイルの内容例:

```env
GEMINI_API_KEY=AIzaSyXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
PORT=8080
```

### 7. 動作確認

```bash
# backend ディレクトリでテストを実行
cd backend
GEMINI_API_KEY=$(grep GEMINI_API_KEY .env | cut -d '=' -f2) go test ./gemini/... -v

# または make コマンドで
make backend-test
```

## トラブルシューティング

### API キーが見つからない

```bash
# すべての API キーを一覧表示
gcloud alpha services api-keys list --project=YOUR_PROJECT_ID
```

### API が有効化されていない

```bash
# 有効化されている API を確認
gcloud services list --enabled

# Generative Language API を有効化
gcloud services enable generativelanguage.googleapis.com
```

### 権限エラー

プロジェクトに対する適切な権限があることを確認してください:

```bash
# 現在のアカウントを確認
gcloud auth list

# プロジェクトのIAMポリシーを確認
gcloud projects get-iam-policy YOUR_PROJECT_ID
```

必要な権限:
- `roles/serviceusage.serviceUsageConsumer` - API の使用
- `roles/apikeys.admin` - API キーの作成・管理

## セキュリティのベストプラクティス

1. **API キーを Git にコミットしない**
   - `.env` ファイルは `.gitignore` に追加済み

2. **API キーに制限を設定する**
   - 特定の API のみに制限
   - 必要に応じて IP アドレス制限

3. **定期的にキーをローテーションする**
   ```bash
   # 新しいキーを作成
   gcloud alpha services api-keys create --display-name="Enjo Simulator API Key (New)"

   # 古いキーを削除
   gcloud alpha services api-keys delete OLD_KEY_ID
   ```

4. **本番環境では Secret Manager を使用する**
   ```bash
   # Secret Manager に保存
   echo -n "YOUR_API_KEY" | gcloud secrets create gemini-api-key --data-file=-

   # Secret Manager から取得
   gcloud secrets versions access latest --secret="gemini-api-key"
   ```

---

## Imagen API (Vertex AI)

画像生成機能に使用するImagen APIのセットアップ方法です。

### 前提条件

- Google Cloud SDK (gcloud) がインストールされていること
- Google Cloud プロジェクトが作成されていること
- プロジェクトで課金が有効化されていること

### 手順

#### 1. Vertex AI API の有効化

```bash
# Vertex AI API を有効化
gcloud services enable aiplatform.googleapis.com

# 有効化を確認
gcloud services list --enabled | grep aiplatform
```

#### 2. Application Default Credentials (ADC) の設定

Imagen APIはAPI keyではなく、Application Default Credentialsを使用します。

```bash
# ADCを設定（初回のみ）
gcloud auth application-default login

# 認証情報を確認
gcloud auth application-default print-access-token
```

これにより、`~/.config/gcloud/application_default_credentials.json` に認証情報が保存されます。

#### 3. プロジェクトIDとロケーションの設定

```bash
# 現在のプロジェクトIDを確認
gcloud config get-value project

# ロケーション一覧を確認（us-central1を推奨）
gcloud compute regions list
```

#### 4. 環境変数の設定

`backend/.env` ファイルに以下を追加:

```env
# GCP Configuration for Imagen
GCP_PROJECT_ID=your-project-id
GCP_LOCATION=us-central1
```

#### 5. 動作確認

```bash
# バックエンドディレクトリに移動
cd backend

# 統合テストを実行（実際のAPI呼び出しが発生するため課金に注意）
export RUN_INTEGRATION_TESTS=true
go test ./image/... -v -run Integration

# 通常のテスト（モックのみ、課金なし）
make backend-test
```

### IAM権限

プロジェクトのサービスアカウントまたはユーザーアカウントに以下の権限が必要:

- `roles/aiplatform.user` - Vertex AI の使用
- または `roles/owner` / `roles/editor` （開発環境の場合）

権限を確認:

```bash
# 現在のアカウントを確認
gcloud auth list

# プロジェクトのIAMポリシーを確認
gcloud projects get-iam-policy YOUR_PROJECT_ID
```

### トラブルシューティング

#### エラー: "Permission denied"

```bash
# ADCを再設定
gcloud auth application-default login

# プロジェクトが正しく設定されているか確認
gcloud config get-value project
```

#### エラー: "API not enabled"

```bash
# Vertex AI APIを有効化
gcloud services enable aiplatform.googleapis.com

# 有効化を確認
gcloud services list --enabled | grep aiplatform
```

#### エラー: "Quota exceeded" または "Rate limit exceeded"

- Imagen APIには生成回数の制限があります
- [GCP Console - Quotas](https://console.cloud.google.com/iam-admin/quotas) で現在の使用状況を確認
- 必要に応じて上限引き上げをリクエスト

### コスト管理

Imagen APIは画像生成ごとに課金されます。コスト管理のため:

```bash
# 予算アラートを設定（GCP Consoleで推奨）
# Billing > Budgets & alerts

# 現在の使用状況を確認
gcloud billing accounts list
```

詳細は [FEATURE_IMAGE_GENERATION.md](./FEATURE_IMAGE_GENERATION.md#コスト見積もり) を参照してください。

---

## Twitter API

Twitter投稿機能に使用するTwitter APIのセットアップ方法です。

**注意**: 現在の実装はモックです。実際のTwitter API統合は将来の拡張として予定されています。

### 将来の実装のための準備

#### 1. Twitter Developer Portal でアプリ作成

1. <https://developer.twitter.com/en/portal/dashboard> にアクセス
2. 「Create App」をクリック
3. アプリ名、説明を入力
4. App Permissions: 「Read and Write」を選択
5. API Key & Secret を取得
6. Access Token & Secret を生成

#### 2. 環境変数の設定

`backend/.env` に追加（将来の実装用）:

```env
# Twitter API Configuration (未実装)
TWITTER_API_KEY=your_twitter_api_key_here
TWITTER_API_SECRET=your_twitter_api_secret_here
TWITTER_ACCESS_TOKEN=your_access_token_here
TWITTER_ACCESS_TOKEN_SECRET=your_access_token_secret_here
```

詳細は [FEATURE_TWITTER_POST.md](./FEATURE_TWITTER_POST.md) を参照してください。

---

## 参考リンク

### Gemini API

- [Google Cloud SDK ドキュメント](https://cloud.google.com/sdk/docs)
- [Generative AI API ドキュメント](https://cloud.google.com/vertex-ai/docs/generative-ai/start/quickstarts/api-quickstart)
- [API キー管理](https://cloud.google.com/docs/authentication/api-keys)

### Imagen API

- [Vertex AI - Imagen Documentation](https://cloud.google.com/vertex-ai/docs/generative-ai/image/overview)
- [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials)
- [Vertex AI Pricing](https://cloud.google.com/vertex-ai/pricing)

### Twitter API

- [Twitter API v2 Documentation](https://developer.twitter.com/en/docs/twitter-api)
- [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard)
