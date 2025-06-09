# ai-aws-extend-switch-roles

AWS Extend Switch Rolesや~/.aws/configに記載するプロファイルを出力するCLIツール

## 概要

このツールは、AWSのアカウント間でのロール切り替えを簡単にするため、AWS Extend Switch Roles 拡張機能や `~/.aws/config`ファイルで使用するプロファイル設定を自動生成します。

## 機能

- AWS Extend Switch Roles用の設定を生成
- `~/.aws/config`用のプロファイル設定を生成

## インストール

### Homebrew

```bash
brew install juliar13/ai-aws-profiles/aws-profiles
```

## 使用方法

### 基本的な使い方

```bash
# 現在のディレクトリに AWS Extend Switch Roles用の設定（extension.ini) と ~/.aws/config 用の設定(config.ini) をそれぞれ出力
aws-prof

# AWS Extend Switch Roles 用の設定を標準出力
aws-prof --format=extension

# ~/.aws/config 用の設定を標準出力
aws-prof --format=config

# 現在のディレクトリの profiles.txt に設定を出力
aws-prof --format=extension --output=profiles.txt

# バージョンの確認
aws-prof --version

# ヘルプ表示
aws-prof --help
```

## 出力例

### AWS Extend Switch Roles

```ini
[profile test-admin]
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole
color = 00aa00

[profile test-readonly]
role_arn = arn:aws:iam::123456789012:role/ReadOnlySwitchRole
color = 00aa00

[profile test2-admin]
role_arn = arn:aws:iam::123456789013:role/AdminSwitchRole
color = 00aa00

[profile test2-readonly]
role_arn = arn:aws:iam::123456789013:role/ReadOnlySwitchRole
color = 00aa00
```

### ~/.aws/config

```ini
[default]
region = region
output = json
role_session_name = user_name

[profile test-admin]
source_profile = default
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole

[profile test-readonly]
source_profile = default
role_arn = arn:aws:iam::123456789012:role/ReadOnlySwitchRole

[profile test2-admin]
source_profile = default
role_arn = arn:aws:iam::123456789013:role/AdminSwitchRole

[profile test2-readonly]
source_profile = default
role_arn = arn:aws:iam::123456789013:role/ReadOnlySwitchRole
```

## オプション

| オプション | 説明 |
|-----------|------|
| `--format` | 出力形式 (extension/config) |
| `--output` | 出力ファイルパス |

## 開発

### 必要な環境

- Go 1.21以上

### ビルド

```bash
go build -o ai-aws-extend-switch-roles cmd/main.go
```

### テスト

```bash
go test ./...
```

## ライセンス

MIT License

## Contributing

Issues やPull Requestsをお待ちしています。
