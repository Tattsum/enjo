# Gemini API キー取得方法

このドキュメントでは、gcloud コマンドを使用して Gemini API キーを取得し、環境変数として設定する方法を説明します。

## 前提条件

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

## 参考リンク

- [Google Cloud SDK ドキュメント](https://cloud.google.com/sdk/docs)
- [Generative AI API ドキュメント](https://cloud.google.com/vertex-ai/docs/generative-ai/start/quickstarts/api-quickstart)
- [API キー管理](https://cloud.google.com/docs/authentication/api-keys)
