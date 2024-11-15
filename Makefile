# makeを打った時のコマンド
.DEFAULT_GOAL := help

.PHONY: install
install: ## 必要なツールをインストール
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: mysql-login
mysql-login: ## dbに入る
	mysql -u root -p -h 127.0.0.1 -P 3307 sns_archive_jwt

# ref: https://zenn.dev/farstep/books/f74e6b76ea7456/viewer/4cd440
.PHONY: create-mg-file
create-mg-file: ## マイグレーションファイルを作成 file_name=xxx
	migrate create -ext sql -dir db/migrations -seq $(file_name)

.PHONY: migrate-up
migrate-up: ## マイグレーションを実行
	migrate -path db/migrations -database "mysql://root:root@tcp(127.0.0.1:3307)/sns_archive_jwt" up

.PHONY: create-down
migrate-down: ## マイグレーションを戻す
	migrate -path db/migrations -database "mysql://root:root@tcp(127.0.0.1:3307)/sns_archive_jwt" down

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
