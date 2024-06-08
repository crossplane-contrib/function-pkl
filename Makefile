REPO := ghcr.io/avarei
IMAGE := function-pkl
TAG := v0.0.0-dev13

.PHONY: build-pkl-crossplane
build-pkl-crossplane:
	pkl project resolve ./pkl/crossplane/
	pkl project package ./pkl/crossplane/

.PHONY: build-pkl-crossplane-example
build-pkl-crossplane-example:
	pkl project resolve ./pkl/crossplane-example/
	pkl project package ./pkl/crossplane-example/

.PHONY: build-image
build-image:
	docker build -t runtime .
	crossplane xpkg build -f package --embed-runtime-image=runtime -o .out/function-pkl.xpkg


.PHONY: push-image
push-image:
	crossplane xpkg push -f .out/function-pkl.xpkg ${REPO}/${IMAGE}:${TAG}
