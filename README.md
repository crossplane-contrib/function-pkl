# function-pkl
[![CI](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml/badge.svg)](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml)

> [!CAUTION]
> This is only a Proof of Concept! This is not yet ready for Prod.

## Pkl Function Call Order
The Composition function is triggered by the Pkl files referenced within the composition pipeline step.
Each Pkl file will be parsed individually. (TODO: add paralellisation?)
* Pkl file within composition
    * The file will amend `CrossplaneResource.pkl`
    * `convert.pkl` is usually called by `import crossplane:state` as a new process.
        * this loads other pkl files defined in Composition CRDs `import crossplane:crds`
        * it also imports pkl-pantry/k8s.contrib/convert.pkl
        * imports pkl-k8s/k8s/k8sResource.pkl
    * `CompositionInput.pkl` is used by the result of `local state = (import crossplane:state) as CompositionInput`

### Basic Pkl File
The bare pkl file we expect is
```pkl
amends "https://raw.githubusercontent.com/Avarei/function-pkl/main/pkl/CrossplaneResource.pkl" // TODO replace with package reference
import "https://raw.githubusercontent.com/Avarei/function-pkl/main/pkl/CompositionInput.pkl" // TODO replace with package reference

local state = import("crossplane:state") as CompositionInput
```

## What does it do?
This Composition function for [Crossplane][crossplane] allows the usage of the [Pkl][pkl] Configuration Language within Compositions.

---

This template uses [Go][go], [Docker][docker], and the [Crossplane CLI][cli] to
build functions.

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime
```

## Debugging
`crossplane beta render example/xr.yaml example/composition.yaml example/functions.yaml --verbose`

### Pkl yaml <-> pkl examples:
Turn Yaml Manifest into Pkl File
```bash
pkl eval -p input=appproject.yaml -o appproject.pkl package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl
pkl eval -p input=example-crd.yaml -o example-crd.pkl package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl
```

Turn Yaml CRD into Pkl template
```bash
pkl eval package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib.crd@1.0.0#/generate.pkl -m . -p source="https://raw.githubusercontent.com/argoproj/argo-cd/master/manifests/crds/appproject-crd.yaml"
```

comparison:
```bash
pkl eval -p input=example-crd.yaml -o example-crd.pkl ../packages/k8s.contrib/convert.pkl
pkl eval -m . -p source="example-crd.yaml" ../packages/k8s.contrib.crd/generate.pkl
```


[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[crossplane]: https://www.crossplane.io
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli
[pkl]: https://pkl-lang.org
