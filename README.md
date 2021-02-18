"# UserService
Goで構築するマイクロサービスのユーザサービス

# 概要
- 簡易的な記事投稿のCRUD機能

## 使用技術
- Go 1.12.17
- Docker docker-compose
- dockerize v0.6.1
- protoc 3.11.0
- gRPC v1.31.0
- AWS(VPC,ECS,ECR,RDS,ELB)
- terraform
- CircleCI

## 構成図
![AWS_stracture](https://user-images.githubusercontent.com/36359899/89097162-79bd3200-d417-11ea-83e5-8c998c824a0f.png)


## 構成図
![PortfolioArchitecture](https://user-images.githubusercontent.com/36359899/108287470-2ef5e280-71ce-11eb-9301-a2c3c8ed5d01.png)

## 機能一覧
- ユーザー
  - 新規登録、編集、削除、全件取得
  - 管理者権限ユーザーによるログイン
  - go-playground/validatorを用いたバリデーション
  - ログインIDとパスワードによる認証・jwt発行
- サービス間通信
  - Envoyプロキシを介した他サービスとの通信

## 関連レポジトリ
- [フロントエンド](https://github.com/yzmw1213/Front)
- [Envoyプロキシ](https://github.com/yzmw1213/Proxy)
- [投稿サービス](https://github.com/yzmw1213/PostService)
