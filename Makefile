
.PHONY: build-pkl-crossplane
build-pkl-crossplane:
	pkl project resolve ./pkl/crossplane/
	pkl project package ./pkl/crossplane/

.PHONY: build-pkl-crossplane-example
build-pkl-crossplane-example:
	pkl project resolve ./pkl/crossplane-example/
	pkl project package ./pkl/crossplane-example/

