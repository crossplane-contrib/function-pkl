REPO ?= github.com/crossplane-contrib/function-pkl
CONTAINER_IMAGE ?= ghcr.io/crossplane-contrib/function-pkl

# Target used for Pkl Package Releases
TARGET =? $(shell git branch --show-current)

LATEST_CORE    := $(shell git tag -l "crossplane.contrib@*.*.*" --sort=-v:refname | head -n 1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
LATEST_EXAMPLE := $(shell git tag -l "crossplane.contrib.example@*.*.*" --sort=-v:refname | head -n 1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
LATEST_XRD     := $(shell git tag -l "crossplane.contrib.xrd@*.*.*" --sort=-v:refname | head -n 1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')

# This Resolves the Dependencies and sets the versions of our packages to the Latest ones for the package in Git
.PHONY: pkl-resolve
pkl-resolve:
	pkl project resolve \
		-e REPOSITORY="$(REPO)" \
 		-e CROSSPLANE_CONTRIB_VERSION="$(LATEST_CORE)" \
		-e CROSSPLANE_CONTRIB_EXAMPLE_VERSION="$(LATEST_EXAMPLE)" \
		-e CROSSPLANE_CONTRIB_XRD_VERSION="$(LATEST_XRD)" \
 		./pkl/*/

.PHONY: pkl-package
pkl-package: pkl-resolve
	$(eval PACKAGE_FILES  := $(shell \
    		pkl project package \
    		 	-e REPOSITORY="$(REPO)" \
    			-e CROSSPLANE_CONTRIB_VERSION="$(LATEST_CORE)" \
    			-e CROSSPLANE_CONTRIB_EXAMPLE_VERSION="$(LATEST_EXAMPLE)" \
    			-e CROSSPLANE_CONTRIB_XRD_VERSION="$(LATEST_XRD)" \
    		 ./pkl/*/ ))

# Ensures the TAG Variable is set.
.PHONY: check-tag
check-tag:
	@[ "${TAG}" ] || (echo "TAG is not specified" && exit 1)

# Initializes Empty Array
RELEASE_FILES :=

# Packages all Projects with the latest tags for each before Pushing the one referenced in TAG
.PHONY: pkl-release
pkl-release: check-tag pkl-package
	$(foreach file,$(PACKAGE_FILES), \
		$(if $(findstring ${TAG},$(file)), \
			$(eval RELEASE_FILES += $(file))))
	@if [ -z "$(RELEASE_FILES)" ]; then \
		echo "No release files found for tag ${TAG}."; \
		exit 1; \
	fi

	gh release create ${TAG} \
	-t ${TAG} \
	-n "" \
	--target ${TARGET} \
	--prerelease \
	--draft \
	$(RELEASE_FILES)

.PHONY: build-image
build-image:
	docker build --build-arg -t runtime .
	crossplane xpkg build -f package --embed-runtime-image=runtime -o .out/function-pkl.xpkg

.PHONY: push-image
push-image:
	crossplane xpkg push -f .out/function-pkl.xpkg ${CONTAINER_IMAGE}:${TAG}
