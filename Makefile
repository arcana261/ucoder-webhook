# Project source files
SRCS = $(patsubst ./%,%,$(shell find . -name "*.go" -not -path "*vendor*"))

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Clean
	rm -f ucoder-webhook

build: ucoder-webhook ## Build

deploy: build ## Deploy service to production (read README.md first)
	ssh www.ucoder.ir "rm -f ~/ucoder-webhook"
	scp ucoder-webhook www.ucoder.ir:~
	ssh -t www.ucoder.ir "sudo bash -c 'mv -f /home/arcana/ucoder-webhook /usr/local/bin/ && systemctl daemon-reload && systemctl enable ucoder-webhook.service && systemctl restart ucoder-webhook.service && systemctl status ucoder-webhook.service'"

ucoder-webhook: $(SRCS)
	go build -o ucoder-webhook

.PHONY: help clean build deploy
