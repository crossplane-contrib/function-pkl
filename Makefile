REPO := ghcr.io/avarei
IMAGE := function-pkl
TAG := v0.0.0-dev13

PKL_MODULE_VERSION_CROSSPLANE := 0.0.10
PKL_MODULE_VERSION_CROSSPLANE_EXAMPLE := 0.0.7

.PHONY: release-pkl-crossplane
release-pkl-crossplane:
	MODULE_REF := crossplane@${PKL_MODULE_VERSION_CROSSPLANE}
	PKL_MODULE_PATH := .out/${MODULE_REF}/${MODULE_REF}
	pkl project resolve ./pkl/crossplane/
	pkl project package ./pkl/crossplane/

	gh release create ${PKL_MODULE_VERSION_CROSSPLANE} \
	-t "${MODULE_REF}" \
	-n "" \
	${PKL_MODULE_PATH} ${PKL_MODULE_PATH}.sha256 ${PKL_MODULE_PATH}.zip ${PKL_MODULE_PATH}.zip.sha256

.PHONY: build-pkl-crossplane-example
release-pkl-crossplane-example:
	MODULE_REF := crossplane-example@${PKL_MODULE_VERSION_CROSSPLANE}
	PKL_MODULE_PATH := .out/${MODULE_REF}/${MODULE_REF}
	pkl project resolve ./pkl/crossplane-example/
	pkl project package ./pkl/crossplane-example/
	gh release create ${PKL_MODULE_VERSION_CROSSPLANE} \
	-t "${MODULE_REF}" \
	-n "" \
	${PKL_MODULE_PATH} ${PKL_MODULE_PATH}.sha256 ${PKL_MODULE_PATH}.zip ${PKL_MODULE_PATH}.zip.sha256

.PHONY: build-image
build-image:
	docker build -t runtime .
	crossplane xpkg build -f package --embed-runtime-image=runtime -o .out/function-pkl.xpkg


.PHONY: push-image
push-image:
	crossplane xpkg push -f .out/function-pkl.xpkg ${REPO}/${IMAGE}:${TAG}
