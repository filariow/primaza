##@ Acceptance Tests

TEST_ACCEPTANCE_OUTPUT_DIR ?= $(OUTPUT_DIR)/acceptance-tests
TEST_ACCEPTANCE_CLI ?= kubectl
TEST_ACCEPTANCE_PARALLEL ?= 4
TEST_ACCEPTANCE_TAGS ?=

ifdef TEST_ACCEPTANCE_TAGS
TEST_ACCEPTANCE_TAGS_ARG ?= --tags="~@disabled" --tags="$(TEST_ACCEPTANCE_TAGS)"
else
TEST_ACCEPTANCE_TAGS_ARG ?= --tags="~@disabled"
endif

.PHONY: setup-venv
setup-venv: ## Setup virtual environment
	python3 -m venv $(PYTHON_VENV_DIR)
	$(PYTHON_VENV_DIR)/bin/pip install --upgrade setuptools
	$(PYTHON_VENV_DIR)/bin/pip install --upgrade pip

.PHONY: test-acceptance-setup
test-acceptance-setup: setup-venv ## Setup the environment for the acceptance tests
	$(PYTHON_VENV_DIR)/bin/pip install -q -r test/acceptance/features/requirements.txt

.PHONY: test-acceptance
test-acceptance: test-acceptance-setup ## Runs acceptance tests
	@(kind get clusters | grep primaza | xargs -I@ kind delete cluster --name @) || true
	echo "Running acceptance tests"
	$(PYTHON_VENV_DIR)/bin/behave --junit --junit-directory $(TEST_ACCEPTANCE_OUTPUT_DIR) --no-capture --no-capture-stderr $(TEST_ACCEPTANCE_TAGS_ARG) $(EXTRA_BEHAVE_ARGS) test/acceptance/features

.PHONY: test-acceptance-x
test-acceptance-x: test-acceptance-setup ## Runs acceptance tests for WIP tagged scenarios
	@(kind get clusters | grep primaza | xargs -I@ kind delete cluster --name @) || true
	echo "Running acceptance tests in $(TEST_ACCEPTANCE_PARALLEL) parallel processes"
	FEATURES_PATH=test/acceptance/features $(PYTHON_VENV_DIR)/bin/behavex -o $(TEST_ACCEPTANCE_OUTPUT_DIR) --no-capture --no-capture-stderr $(TEST_ACCEPTANCE_TAGS_ARG) $(EXTRA_BEHAVE_ARGS) --parallel-processes $(TEST_ACCEPTANCE_PARALLEL)

.PHONY: clean
clean: ## Removes temp directories
	-rm -rf ${V_FLAG} $(OUTPUT_DIR)