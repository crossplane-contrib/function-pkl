# Deploying this Example

## Prerequirements
* [pkl cli](https://pkl-lang.org/main/current/pkl-cli/index.html#installation)
* crossplane in a cluster
* function-pkl deployed (see [here](../../README.md))
* provider-kubernetes deployed
* provider config for kubernetes provider called "default"

## Deploy XRD
```shell
cd pkl/pkl/crossplane.contrib.example
pkl eval xrds/ExampleXR.pkl | kubectl apply -f -
```
## Deploy Composition
```shell
kubectl apply -f ../../example/full/composition.yaml
```
## Deploy XR
```shell
kubectl apply -f ../../example/full/xr.yaml
```

## Check the Resource
```shell
kubectl get xrs.example.crossplane.io example-xr -oyaml

crossplane beta trace xrs example-xr
```
