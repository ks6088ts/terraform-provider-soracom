TERRAFORM_DIR ?= $(PWD)
TERRAFORM_DIR_LIST ?= $(wildcard examples/*)
TERRAFORM_PLUGIN_DOCS_VERSION ?= 0.13.0
TERRAFORM ?= cd $(TERRAFORM_DIR) && terraform
GENERATED_DIR ?= docs

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.DEFAULT_GOAL := help

.PHONY: install-tfplugindocs
install-tfplugindocs:
	which tfplugindocs || go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v$(TERRAFORM_PLUGIN_DOCS_VERSION)

.PHONY: docs-generate
docs-generate:
	tfplugindocs

.PHONY: docs-diff
docs-diff:
	git diff --exit-code --relative $(GENERATED_DIR)

.PHONY: ci-test-docs
ci-test-docs: install-tfplugindocs docs-generate ## ci test for documents (fixme: run `docs-diff`)

.PHONY: clear
clear:
	cd $(TERRAFORM_DIR) && rm -rf .terraform*

.PHONY: init
init:
	$(TERRAFORM) init

.PHONY: lint
lint:
	$(TERRAFORM) validate

.PHONY: format-check
format-check:
	$(TERRAFORM) fmt -check

.PHONY: format
format: ## format terraform codes
	$(TERRAFORM) fmt -recursive

.PHONY: plan
plan:
	$(TERRAFORM) plan -var-file=$(PWD)/test/test.tfvars

.PHONY: _ci-test-base
_ci-test-base: clear init lint format-check plan

.PHONY: ci-test-examples
ci-test-examples: ## ci test for examples
	for dir in $(TERRAFORM_DIR_LIST) ; do \
		make -f terraform.mk _ci-test-base TERRAFORM_DIR=$$dir || exit 1 ; \
	done

.PHONY: ci-test
ci-test: ci-test-examples ci-test-docs ## ci test
