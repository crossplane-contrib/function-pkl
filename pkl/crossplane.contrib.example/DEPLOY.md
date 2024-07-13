# Deploying this Example

## Prerequirements
* [pkl cli](https://pkl-lang.org/main/current/pkl-cli/index.html#installation) and kubectl are installed
* a kubernetes cluster with crossplane deployed
* function-pkl deployed (see [README.md](../../README.md))
* [provider-kubernetes](https://marketplace.upbound.io/providers/crossplane-contrib/provider-kubernetes) deployed
* a [provider config](https://marketplace.upbound.io/providers/crossplane-contrib/provider-kubernetes/v0.14.0/resources/kubernetes.crossplane.io/ProviderConfig/v1alpha1) for provider-kubernetes called "default"

## Deploy XRD
```shell
pkl eval package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane.contrib.example@0.0.1#/xrds/ExampleXR.pkl | kubectl apply -f -
```
## Deploy Composition
```shell
pkl eval package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane.contrib.example@0.0.1#/compositions/uri.pkl | kubectl apply -f -
```
## Deploy XR
```shell
pkl eval package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane.contrib.example@0.0.1#/xrs/uri.pkl | kubectl apply -f -
```

## Check the Resource
```shell
kubectl get xrs.example.crossplane.io uri-example -oyaml
```

```shell
crossplane beta trace xrs uri-example
```
