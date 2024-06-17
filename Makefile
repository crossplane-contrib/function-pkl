REPO := github.com/crossplane-contrib/function-pkl
CONTAINER_IMAGE := ghcr.io/crossplane-contrib/function-pkl
TAG := v0.0.1

# Branch used for Pkl Package Releases
BRANCH := $(shell git branch --show-current)

PKL_BASE_URI := package://pkg.pkl-lang.org

PKL_CORE_NAME := crossplane
PKL_CORE_VERSION := 0.0.29
PKL_CORE_REF := ${PKL_CORE_NAME}@${PKL_CORE_VERSION}
PKL_CORE_URI := ${PKL_BASE_URI}/${REPO}/${PKL_CORE_REF}

PKL_EXAMPLE_NAME := crossplane-example
PKL_EXAMPLE_VERSION := 0.1.19
PKL_EXAMPLE_REF := ${PKL_EXAMPLE_NAME}@${PKL_EXAMPLE_VERSION}
PKL_EXAMPLE_URI := ${PKL_BASE_URI}/${REPO}/${PKL_CORE_REF}

SED_MAC = sed -i '' -E "s|($(PKL_BASE_URI)/${REPO}/$(PACKAGE_NAME)@)([0-9]+\.[0-9]+\.[0-9]+)|\1$(PACKAGE_VERSION)|g"
SED_LINUX = sed -i -E "s|($(PKL_BASE_URI)/${REPO}/$(PACKAGE_NAME)@)([0-9]+\.[0-9]+\.[0-9]+)|\1$(PACKAGE_VERSION)|g"
SED_TARGETS := example/ README.md pkl/${PKL_CORE_NAME}/PklProject pkl/${PKL_EXAMPLE_NAME}/PklProject

.PHONY: build-core-package
build-core-package: PACKAGE_NAME := $(PKL_CORE_NAME)
build-core-package: PACKAGE_VERSION := $(PKL_CORE_VERSION)
build-core-package:
	pkl project resolve ./pkl/${PACKAGE_NAME}/
ifeq ($(shell uname), Darwin)
	find $(SED_TARGETS) -type f -exec $(SED_MAC) {} +
else
	find $(SED_TARGETS) -type f -exec $(SED_LINUX) {} +
endif


.PHONY: release-pkl-crossplane
release-pkl-crossplane:
	pkl project resolve ./pkl/${PKL_CORE_NAME}/ && \
	RELEASE_FILES=$$(pkl project package ./pkl/${PKL_CORE_NAME}/) && \
	gh release create ${PKL_CORE_REF} \
	-t ${PKL_CORE_REF} \
	-n "" \
	--target ${BRANCH} \
	--prerelease \
	$$RELEASE_FILES

.PHONY: build-example-package
build-example-package: PACKAGE_NAME := $(PKL_EXAMPLE_NAME)
build-example-package: PACKAGE_VERSION := $(PKL_EXAMPLE_VERSION)
build-example-package:
	pkl project resolve ./pkl/${PACKAGE_NAME}/
ifeq ($(shell uname), Darwin)
	find $(SED_TARGETS) -type f -exec $(SED_MAC) {} +
else
	find $(SED_TARGETS) -type f -exec $(SED_LINUX) {} +
endif

.PHONY: build-pkl-crossplane-example
release-pkl-crossplane-example:
	pkl project resolve ./pkl/${PKL_EXAMPLE_NAME}/ && \
	RELEASE_FILES=$$(pkl project package ./pkl/${PKL_EXAMPLE_NAME}/) && \
	gh release create ${PKL_EXAMPLE_REF} \
	-t ${PKL_EXAMPLE_REF} \
	-n "" \
	--target ${BRANCH} \
	--prerelease \
	$$RELEASE_FILES

.PHONY: build-image
build-image:
	docker build --build-arg PKL_CORE_PACKAGE=${PKL_CORE_URI} -t runtime .
	crossplane xpkg build -f package --embed-runtime-image=runtime -o .out/function-pkl.xpkg


.PHONY: push-image
push-image:
	crossplane xpkg push -f .out/function-pkl.xpkg ${CONTAINER_IMAGE}:${TAG}
