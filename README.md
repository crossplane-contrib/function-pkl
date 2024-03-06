# function-pkl
[![CI](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml/badge.svg)](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml)

For now this is only a Proof of Concept.

For the moment this function creates the following filetree in each invocation:
```yaml
PklProject: <optional. may contain dependencies for using the @notation>
Example.pkl: # name of the CompositeResourceDefinition
  # contains the pkl template of the CompositeResourceDefinition.
  # used to transform the yaml representation of the CompositeResource to the Pkl file
  # - desired.composition
  # - observed.composition
observed:
  composition: <(XR) contains the pkl file representation of the Composite Resource>
  resources:
    aPod.pkl: <A Pod Manifest containing status fields>
    # aConfigMap.pkl has not yet been created in this example.
desired:
  composition: <amends "/Example.pkl">
  resources:
    aPod.pkl: <A Pod Manifest in Pkl File format>
    aConfigMap.pkl: <A ConfigMap which could e.g. import "/observed/aPod.pkl" to read it's status>
```

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

# Create Own Reader

## Kinds of Readers needed
- crossplane:  
  the reader will provide access to everything in observed and desired.

### Deep dive
As a crossplane Function I get observed, containing composition, containing resource, containing 

## Implement it like this:
[example Implementation](https://github.com/apple/pkl-go/blob/main/cmd/internal/gen-fixtures/gen-fixtures.go#L119)

## Implement the following Interface:
- 



[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli
