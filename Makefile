# We want to have our binaries in the bin subdirectory available. In addition we want them to have priority over
# binaries somewhere else on the system.
export PATH := $(shell pwd)/bin:$(PATH)

.PHONY: all
all: build

V1ALPHA1_DEEPCOPY_FILE := api/v1alpha1/zz_generated.deepcopy.go
V1ALPHA1_TYPE_FILES := $(filter-out $(V1ALPHA1_DEEPCOPY_FILE), $(wildcard api/v1alpha1/*.go))
$(V1ALPHA1_DEEPCOPY_FILE): $(V1ALPHA1_TYPE_FILES)
	controller-gen object paths=./api/v1alpha1/...

V1ALPHA1_CRD_FILES := \
	core.ctf.backbone81_apikeys.yaml \
	core.ctf.backbone81_challengedescriptions.yaml \
	core.ctf.backbone81_challengeinstances.yaml
V1ALPHA1_CRD_FILES := $(addprefix config/crd/bases/,$(V1ALPHA1_CRD_FILES))
$(V1ALPHA1_CRD_FILES): $(V1ALPHA1_TYPE_FILES)
	controller-gen crd paths=./api/v1alpha1/... output:crd:artifacts:config=config/crd/bases

V1ALPHA1_ROLE_FILES := \
	role.yaml
V1ALPHA1_ROLE_FILES := $(addprefix config/rbac/,$(V1ALPHA1_ROLE_FILES))
V1ALPHA1_CONTROLLER_FILES := $(shell go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}}{{"\n"}}{{end}}' ./internal/controller/...)
$(V1ALPHA1_ROLE_FILES): $(V1ALPHA1_CONTROLLER_FILES)
	controller-gen rbac:roleName=ctf-challenge-operator paths=./internal/controller/...

.PHONY: generate
generate: $(V1ALPHA1_DEEPCOPY_FILE) $(V1ALPHA1_CRD_FILES) $(V1ALPHA1_ROLE_FILES) ## Generate files

.PHONY: prepare
prepare: generate ## Run go fmt against code.
	go mod tidy
	go fmt ./...
	go vet ./...

.PHONY: lint
lint: prepare ## Run linter
	golangci-lint run --fix

.PHONY: build
build: lint ## Build manager binary.
	go build ./cmd/ctf-challenge-operator

.PHONY: run
run: lint ## Run a controller from your host.
	go run ./cmd/ctf-challenge-operator --enable-developer-mode --log-level 10

# ========== legacy makefile starting from here ==========

# Image URL to use all building/pushing image targets
IMG ?= controller:latest

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: test
test: lint ## Run tests.
	go test $$(go list ./... | grep -v /e2e) -coverprofile cover.out

# TODO(user): To use a different vendor for e2e tests, modify the setup under 'tests/e2e'.
# The default setup assumes Kind is pre-installed and builds/loads the Manager Docker image locally.
# CertManager is installed by default; skip with:
# - CERT_MANAGER_INSTALL_SKIP=true
.PHONY: test-e2e
test-e2e: lint ## Run the e2e tests. Expected an isolated environment using Kind.
	@command -v kind >/dev/null 2>&1 || { \
		echo "Kind is not installed. Please install Kind manually."; \
		exit 1; \
	}
	@kind get clusters | grep -q 'kind' || { \
		echo "No Kind cluster is running. Please start a Kind cluster before running the e2e tests."; \
		exit 1; \
	}
	go test ./test/e2e/ -v -ginkgo.v

##@ Build

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker push ${IMG}

.PHONY: build-installer
build-installer: generate ## Generate a consolidated YAML with CRDs and deployment.
	mkdir -p dist
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default > dist/install.yaml

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install
install: generate ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kustomize build config/crd | kubectl apply -f -

.PHONY: uninstall
uninstall: generate ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/crd | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

.PHONY: deploy
deploy: generate ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && kustomize edit set image controller=${IMG}
	kustomize build config/default | kubectl apply -f -

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster specified in ~/.kube/config. Call with ignore-not-found=true to ignore resource not found errors during deletion.
	kustomize build config/default | kubectl delete --ignore-not-found=$(ignore-not-found) -f -
