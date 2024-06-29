# function-pkl
[![CI](https://github.com/crossplane-contrib/function-pkl/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane-contrib/function-pkl/actions/workflows/ci.yml)

This [composition function][functions] allows to deploy [Crossplane Managed Resources][crossplane] from [Pkl Packages][pkl].

## Installation
```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: Function
metadata:
  name: function-pkl
spec:
  package: ghcr.io/crossplane-contrib/function-pkl:v0.0.1
```

## Usage
To use this function create add the following to your composition:
```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: function-pkl
spec:
  compositeTypeRef:
    apiVersion: example.crossplane.io/v1
    kind: XR
  mode: Pipeline
  pipeline:
  - step: run-the-template
    functionRef:
      name: function-pkl
    input:
      apiVersion: template.fn.crossplane.io/v1beta1
      kind: Pkl
      spec:
        type: uri
        # This pkl file is at `pkl/crossplane.contrib.example/full.pkl` in this repo
        uri: "package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane.contrib.example@1.0.0#/full.pkl"
```

### Example
see [examples](./example/)
and see [pkl/crossplane.contrib.example/DEPLOY.md](pkl/crossplane.contrib.example/DEPLOY.md) on instructions how to deploy it.
## Creating a new Composition Function
see [pkl/crossplane.contrib.example/DEVELOP.md](pkl/crossplane.contrib.example/DEVELOP.md)

## Development
This function uses [Go][go], [Docker][docker], the [Crossplane CLI][cli], and the [Pkl CLI][pkl cli] to build functions

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests - see fn_test.go
$ go test ./...

# Build the function's runtime image - see Dockerfile
$ docker build . --tag=runtime

# Build a function package - see package/crossplane.yaml
$ crossplane xpkg build -f package --embed-runtime-image=runtime

# Push the Package to a registry
$ crossplane xpkg push -f function-pkl*.xpkg ghcr.io/crossplane-contrib/function-pkl

# Resolve Pkl Project Dependencies
$ pkl project resolve

# Package the Pkl Project
$ pkl project package

# Debugging this function
# While a Debugging session is running run on the same host:
$ crossplane beta render xr.yaml composition.yaml functions.yaml --verbose
```

### Pkl Function Flow
This Chart illustrates how what happens, when Crossplane Triggers this Composition Function.

```mermaid
sequenceDiagram
    participant fun as RunFunction
    box rgb(20, 100, 60) Pkl Files
        participant full.pkl as full.pkl
        participant c.pkl as Crossplane.pkl
    end
    box rgb(20,60,100) CrossplaneReader
        participant c.r as crossplane:request
    end

    fun->>fun: Create Pkl Evaluator
    activate fun
    fun->>full.pkl: Evaluate Pkl File
    full.pkl->>c.pkl: import
    c.pkl->>full.pkl: Module with helper functions
    full.pkl->>+c.pkl: get CompositionRequest
        c.pkl->>c.r: Get RunFunctionRequest
        c.r->>c.pkl: Return it in yaml
        c.pkl->>c.pkl: convert from yaml to Pkl Objects
    c.pkl->>-full.pkl: return CompositionRequest
    full.pkl->>fun: Function Respone
    deactivate fun
```


[functions]: https://docs.crossplane.io/latest/concepts/composition-functions
[go]: https://go.dev
[crossplane]: https://www.crossplane.io
[function guide]: https://docs.crossplane.io/knowledge-base/guides/write-a-composition-function-in-go
[package docs]: https://pkg.go.dev/github.com/crossplane/function-sdk-go
[docker]: https://www.docker.com
[cli]: https://docs.crossplane.io/latest/cli
[pkl]: https://pkl-lang.org
[pkl cli]: https://pkl-lang.org/main/current/pkl-cli/index.html#installation
[provider-kubernetes]: https://marketplace.upbound.io/providers/crossplane-contrib/provider-kubernetes/v0.13.0
[provider-kubernetes-object]: https://marketplace.upbound.io/providers/crossplane-contrib/provider-kubernetes/v0.13.0/resources/kubernetes.crossplane.io/Object/v1alpha2
