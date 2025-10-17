# プロジェクトルール

## 開発思想 (t-wada style)

このプロジェクトは以下の原則に従います:

1. **テスト駆動開発 (TDD)**: テストファーストで開発する
2. **小さく作って育てる**: 一度に完璧を目指さず、動くものを小さく作り、少しずつ改善する
3. **シンプルさを保つ**: 複雑さは敵。最もシンプルな解決策を選ぶ
4. **継続的なリファクタリング**: テストがあるから安心してリファクタリングできる

## TDD サイクル (Red-Green-Refactor)

```text
1. Red:   失敗するテストを書く
2. Green: テストが通る最小限のコードを書く
3. Refactor: コードをきれいにする
4. Repeat: 次の機能へ
```

## 開発ワークフロー

### Makefileを使った開発

```bash
# ヘルプを表示
make help

# すべてのチェックを実行 (fmt → lint → test)
make check

# バックエンドのみチェック
make backend-check

# フロントエンドのみチェック
make frontend-check
```

### バックエンド (Go)

```bash
# TDDサイクル
cd backend

# 1. テストを書く (Red)
# 2. テストを実行して失敗を確認
go test ./... -v

# 3. 実装を書く (Green)
# 4. テストが通ることを確認
go test ./... -v

# 5. リファクタリング (Refactor)
# 6. フォーマット・Lint・テストを実行
make backend-check
```

### フロントエンド (TypeScript/Next.js)

```bash
# TDDサイクル
cd frontend

# 1. テストを書く (Red)
# 2. テストを実行して失敗を確認
npm run test

# 3. 実装を書く (Green)
# 4. テストが通ることを確認
npm run test

# 5. リファクタリング (Refactor)
# 6. すべてのチェックを実行
make frontend-check
```

## 実装順序 (TDD)

1. **テストを書く** - 期待する動作をテストコードで表現
2. **テストを実行** - 失敗することを確認 (Red)
3. **最小限の実装** - テストが通る最小限のコードを書く (Green)
4. **テストを実行** - 成功することを確認
5. **リファクタリング** - コードをきれいにする (Refactor)
6. **フォーマット・Lint** - `make check` でコード品質を確認
7. **次の機能へ** - 小さいステップで繰り返す

## コミット前チェックリスト

- [ ] コードフォーマット済み
- [ ] Lintエラーなし
- [ ] 全テスト通過
- [ ] 新機能にテストを追加済み

## コーディング規約

### Go
- エラーハンドリングは必須
- コンテキストを適切に使用
- interfaceは小さく保つ
- gofmtに従う

### TypeScript
- 型定義を明示的に記述
- any型の使用を避ける
- React Hooksのルールに従う
- コンポーネントは単一責任の原則に従う

## Docker開発環境

```bash
# 環境構築
docker-compose up --build

# バックエンドのみ再起動
docker-compose restart backend

# フロントエンドのみ再起動
docker-compose restart frontend

# ログ確認
docker-compose logs -f backend
docker-compose logs -f frontend
```
