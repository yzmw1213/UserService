# GoMicroApp
Goで構築するマイクロサービスの雛形アプリケーション

# 概要
- 簡易的な記事投稿のCRUD機能

## 使用技術
- Go 1.12.17
- Docker docker-compose
- protoc 3.11.0
- dockerize v0.6.1
- AWS(VPC,ECS,ECR,MySQL,ELB)
- terraform
- CircleCI

## 構成図
![AWS_stracture](https://user-images.githubusercontent.com/36359899/89095633-d534f300-d40a-11ea-99a9-dd4b14e32920.png)

## 機能一覧
- 簡易投稿
  - 新規登録、編集、削除、全件取得
  - go-playground/validatorを用いたバリデーション