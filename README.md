# function-pkl
[![CI](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml/badge.svg)](https://github.com/Avarei/function-pkl/actions/workflows/ci.yml)

This [composition function][functions] allows to deploy [Crossplane Managed Resources][crossplane] from [Pkl Packages][pkl].

## Usage
Currently the Composite Resource Definition, as well as All CRDs must be available in the Pkl Format, and deployed as a Pkl Package so that they can be referenced as observed variables.

### Example
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
        pklCRDs:
          - apiVersion: example.crossplane.io/v1
            kind: XR
            uri: package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.2#/crds/XR.pkl
          - apiVersion: kubernetes.crossplane.io/v1alpha2
            kind: Object
            uri: package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.2#/crds/Object.pkl
        pklComposition:
          name: XR
          type: uri
          uri: package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.2#/crds/XR.pkl
        pklManifests:
          - name: object-one
            type: uri
            uri: package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.2#/object-one.pkl
          - name: object-two
            type: uri
            uri: package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane-example@0.0.2#/object-two.pkl
```

## Building a Pkl Package
A Pkl Package can be built in the following steps:
1. Create Pkl files in a directory (you can convert CRDs, XRDs and Manifests to Pkl) # TODO Link
1. Create PklProject file (take a look at pkl/crossplane-example/PklProject for reference)
    * Make sure that `\(baseUri)@\(version)` is the url where the package metadata can be downloaded from
    * Likewise `\(baseUri)@\(version).zip` is the path where the content will be expected.
    * The PklProjects in this repository use pkg.pkl-lang.org, as they redirect to the download
        * the path https://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.15
        * redirects to https://github.com/avarei/function-pkl/releases/download/crossplane@0.0.15/crossplane@0.0.15
1. Run `pkl project resolve` to resolve the dependencies of the PklProject
1. Run `pkl project package` to resolve the dependencies of PklProject
1. Publish the Package
    * This can be done by uploading the files created by the previous command to the place defined in step 2.
    * in github this can be done by creating a release and uploading the files to it.

## Basic Pkl File
The bare pkl file we expect is
```pkl
amends "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.15#/CrossplaneResource.pkl"
import "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.15#/CompositionInput.pkl"

local state = import("crossplane:state") as CompositionInput
```

### Pkl Function Call Order
The Composition function is triggered by the Pkl files referenced within the composition pipeline step.
Each Pkl file will be parsed individually. (TODO: add paralellisation?)
* Pkl file within composition
    * The file will amend `CrossplaneResource.pkl`
    * `convert.pkl` is usually called by `import crossplane:state` as a new process.
        * this loads other pkl files defined in Composition CRDs `import crossplane:crds`
        * it also imports pkl-pantry/k8s.contrib/convert.pkl
        * imports pkl-k8s/k8s/k8sResource.pkl
    * `CompositionInput.pkl` is used by the result of `local state = (import crossplane:state) as CompositionInput`

### Generating Pkl Files and Modules from Manifests
There are some package to make it easier to convert existing CRDs or Manifests into the Pkl format.
> [!NOTE]
> There is currently no module to convert from CompositeResourceDefinitions (XRDs) to Pkl.
> Crossplane Creates CRDs from the XRDs. These CRDs can be converted instead.
> If direct conversion would be useful to you - please create an issue.

In this example we'll implement the [Object][provider-kubernetes-object] resource from [provider-kubernetes][provider-kubernetes].

1. First we will generate a Pkl Module from the Object CRD.
   ```bash
   pkl eval "package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib.crd@1.0.4#/generate.pkl" -m . -p source="https://raw.githubusercontent.com/crossplane-contrib/provider-kubernetes/main/package/crds/kubernetes.crossplane.io_objects.yaml"
   ```
   This should create a `Object.pkl` file.

1. Create a small Pkl file which adds context for the converter about the CRDs.
   ```pkl
   amends "package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib@1.0.1#/convert.pkl"
   customResourceTemplates {
     ["Object"] {
       ["kubernetes.crossplane.io/v1alpha2"] = import("Object.pkl")
     }
   }
   ```

1. Convert a object manifest to pkl
   ```bash
   pkl eval -p "input=https://raw.githubusercontent.com/crossplane-contrib/provider-kubernetes/main/examples/object/object.yaml" -o example-object.pkl <CONVERT-FILE-FROM-PREVIOUS-STEP>.pkl
   ```

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
$ crossplane xpkg push -f function-pkl*.xpkg ghcr.io/avarei/function-pkl

# Resolve Pkl Project Dependencies
$ pkl project resolve

# Package the Pkl Project
$ pkl project package

# Debugging this function
# While a Debugging session is running run on the same host:
$ crossplane beta render xr.yaml composition.yaml functions.yaml --verbose
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
