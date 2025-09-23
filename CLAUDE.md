# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

このリポジトリは、作業で使用するツール集です。各ディレクトリは独立したツールとして機能し、主にGoとNode.jsで実装されています。

## プロジェクト構造

### Goプロジェクト
- `discordBotToTranslation/` - Discord翻訳ボット
- `createTestingIssue/` - PR向けテスト用issue作成
- `createTestingIssueFromMilestone/` - マイルストーンからのテスト用issue作成
- `createIssueBySpreadsheets/` - スプレッドシートからのissue作成
- `postAssignableIssues/` - アサイン可能なissue投稿

### Node.jsプロジェクト
- `mfmCsv/` - CSV処理ツール（金融データ抽出用）

### その他
- `postSlackFromSpreadsheet/` - Google Apps Scriptプロジェクト
- `iosShortcut/` - iOS Shortcut用スクリプト

## 開発用コマンド

### Goプロジェクト
各Goプロジェクトは独立して動作します。

```bash
# ローカル実行（例：Discord翻訳ボット）
cd discordBotToTranslation
go run main.go

# テスト実行（もしテストファイルがある場合）
go test ./...

# ビルド
go build
```

### Node.jsプロジェクト（mfmCsv）
```bash
cd mfmCsv
npm install
npm start

# 実行例
node index.js "2025-02-01_2025-02-28.csv" "〇〇銀行"
```

## 環境変数設定

各プロジェクトで必要な環境変数は以下の通りです：

### discordBotToTranslation
```bash
export DISCORD_BOT_TOKEN=xxxxxx
export JAPANESE_CHANNEL_ID=xxxxxx
export ENGLISH_CHANNEL_ID=xxxxxx
export GOOGLE_APPLICATION_CREDENTIALS=xxxxxx
```

### GitHub関連ツール（createTestingIssue系）
```bash
export GITHUB_APP_PRIVATE_KEY=$(cat private-key.pem)
export GITHUB_APP_ID=xxxxxx
export GITHUB_OWNER=xxxxxx
export INSTALLATION_ID=xxxxxx
export SLACK_TOKEN=xxxxxx
export SLACK_CHANNEL=xxxxxx
```

## アーキテクチャ

各ツールは独立しており、モノリポ形式で管理されています：

1. **Goベースのツール**: GitHub APIやSlack API連携が中心
2. **Node.jsベースのツール**: CSV処理や文字エンコーディング変換が中心
3. **Google Apps Script**: スプレッドシートとSlack連携用

## 注意事項

- 各プロジェクトは独立した依存関係を持っています
- Goプロジェクトは1.19以上を想定
- Node.jsプロジェクトはES Modulesを使用
- 認証情報は環境変数で管理してください