# AWS プロファイルの設定出力ツール

AWS Extend Switch Roles や `~/.aws/config` に記載するプロファイルを出力するCLIツール

## 概要

このツールは、AWSのアカウント間でのロール切り替えを簡単にするため、AWS Extend Switch Roles 拡張機能や `~/.aws/config`ファイルで使用するプロファイル設定を自動生成します。

## 機能

- AWS Organizations から自動的にアカウント情報を取得
- プロファイル設定を生成
  - AWS Extend Switch Roles 用の設定を生成
  - `~/.aws/config`用のプロファイル設定を生成

## 前提条件

- AWS CLI がセットアップされていること
- AWS Organizations の管理アカウントにアクセス権限があること
- Organizations の ListAccounts 権限があること

## インストール

### Homebrew

対応予定です。

### ソースからビルド

```bash
git clone https://github.com/juliar13/ai-aws-profiles.git
cd ai-aws-profiles
go build -o aws-prof cmd/main.go
sudo cp aws-prof /usr/local/bin/
```

## 使用方法

### 基本的な使い方

```bash
# AWS Organizations から自動取得して extension.ini と config.ini をそれぞれ出力
aws-prof

# AWS Extend Switch Roles 用の設定を標準出力
aws-prof --format=extension

# ~/.aws/config 用の設定を標準出力
aws-prof --format=config

# 指定したファイルに設定を出力
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
role_arn = arn:aws:iam::123456789011:role/AdminSwitchRole
region = ap-northeast-1
color = 00aa00

[profile test-readonly]
role_arn = arn:aws:iam::123456789011:role/ReadOnlySwitchRole
region = ap-northeast-1
color = 00aa00

[profile test2-admin]
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole
region = ap-northeast-1
color = 00aa00

[profile test2-readonly]
role_arn = arn:aws:iam::123456789012:role/ReadOnlySwitchRole
region = ap-northeast-1
color = 00aa00
```

### ~/.aws/config

```ini
[default]
region = ap-northeast-1
output = json
role_session_name = user_name

[profile test-admin]
source_profile = default
role_arn = arn:aws:iam::123456789011:role/AdminSwitchRole

[profile test-readonly]
source_profile = default
role_arn = arn:aws:iam::123456789011:role/ReadOnlySwitchRole

[profile test2-admin]
source_profile = default
role_arn = arn:aws:iam::123456789012:role/AdminSwitchRole

[profile test2-readonly]
source_profile = default
role_arn = arn:aws:iam::123456789012:role/ReadOnlySwitchRole
```

## オプション

| オプション | 説明 |
|-----------|------|
| `--format` | 出力形式 (extension/config) |
| `--output` | 出力ファイルパス |

## AWS Extend Switch Role の色を設定

color-setting.ini に記載されているように、上から順に反映させること。
以下はデフォルトの設定だが、プロファイル名に admin が含まれているプロファイルは、color = 6644FF となり、 readonly が含まれていれば color = 22CCAA となる。
例えば、プロファイル名に admin と readonly が含まれていると、下にある readonly 22CCAA が適用されて、color = 22CCAA となる。

```
admin 6644FF
readonly 22CCAA
```

```ini
[profile test3-admin]
role_arn = arn:aws:iam::123456789013:role/AdminSwitchRole
region = ap-northeast-1
color = 6644FF

[profile test3-readonly]
role_arn = arn:aws:iam::123456789013:role/ReadOnlySwitchRole
region = ap-northeast-1
color = 22CCAA

[profile test4-admin]
role_arn = arn:aws:iam::123456789014:role/AdminSwitchRole
region = ap-northeast-1
color = 6644FF

[profile test4-readonly]
role_arn = arn:aws:iam::123456789014:role/ReadOnlySwitchRole
region = ap-northeast-1
color = 22CCAA
```


## 開発

### 必要な環境

- Go 1.21以上

### ビルド

```bash
go build -o aws-prof cmd/main.go
```

### テスト

```bash
go test ./...
```

## ライセンス

MIT License

## Contributing

Issues やPull Requestsをお待ちしています。
