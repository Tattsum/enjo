#!/bin/bash
# Gemini API キーの動作確認スクリプト

set -e

# カラー出力
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 Gemini API キーの動作確認を開始します...${NC}\n"

# .env ファイルの存在確認
if [ ! -f "backend/.env" ]; then
    echo -e "${RED}❌ backend/.env ファイルが見つかりません${NC}"
    echo -e "${YELLOW}💡 以下のコマンドで .env ファイルを作成してください:${NC}"
    echo "   cp backend/.env.example backend/.env"
    echo "   # そして GEMINI_API_KEY を設定してください"
    exit 1
fi

# .env から API キーを読み込む
source backend/.env

# API キーの存在確認
if [ -z "$GEMINI_API_KEY" ]; then
    echo -e "${RED}❌ GEMINI_API_KEY が設定されていません${NC}"
    echo -e "${YELLOW}💡 backend/.env ファイルに GEMINI_API_KEY を設定してください${NC}"
    echo ""
    echo "gcloud コマンドで取得する場合:"
    echo "  gcloud alpha services api-keys create --display-name=\"Enjo Simulator API Key\""
    echo "  gcloud alpha services api-keys list"
    echo ""
    echo "詳細は docs/SETUP_API_KEY.md を参照してください"
    exit 1
fi

# API キーの形式確認
if [[ ! $GEMINI_API_KEY =~ ^AIza ]]; then
    echo -e "${YELLOW}⚠️  警告: API キーが正しい形式でない可能性があります${NC}"
    echo "   Gemini API キーは通常 'AIza' で始まります"
    echo ""
fi

echo -e "${GREEN}✅ API キーが設定されています${NC}"
echo -e "   キー: ${GEMINI_API_KEY:0:10}...${GEMINI_API_KEY: -4}\n"

# Go テストの実行
echo -e "${YELLOW}🧪 Gemini API との接続テストを実行します...${NC}\n"

cd backend

# テストの実行
if GEMINI_API_KEY=$GEMINI_API_KEY go test -v -run TestClient_GenerateInflammatoryText ./gemini/...; then
    echo -e "\n${GREEN}✅ API キーは正常に動作しています！${NC}"
    echo -e "${GREEN}🎉 Gemini API との接続に成功しました${NC}\n"
    exit 0
else
    echo -e "\n${RED}❌ API との接続に失敗しました${NC}"
    echo -e "${YELLOW}💡 以下を確認してください:${NC}"
    echo "   1. API キーが正しく設定されているか"
    echo "   2. Generative Language API が有効化されているか"
    echo "      gcloud services enable generativelanguage.googleapis.com"
    echo "   3. API キーに制限が設定されている場合、正しく設定されているか"
    echo ""
    echo "詳細は docs/SETUP_API_KEY.md を参照してください"
    exit 1
fi
