# 🧰 Makefile - Plataforma Cryptos
#
# Este arquivo contém as tarefas para desenvolvimento, build e deploy.
#

.PHONY: install update verify

# ------------------------------------------------------------
# 🛠️ Setup e Desenvolvimento
# ------------------------------------------------------------

install:
	@go mod download
	@echo "✅ Dependências baixadas com sucesso!"

update:
	@go get -u ./...
	@echo "✅ Dependências atualizadas com sucesso!"

verify:
	@rm -rf go.sum
	@go clean -modcache
	@go mod tidy
	@go mod verify
	@echo "✅ Dependências verificadas com sucesso"
