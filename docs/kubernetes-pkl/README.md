# Converting any Kubernetes Manifest to Pkl

## (If using CRDs) Convert Manifest of CRD to Pkl Template
If Manifests to be converted are Custom Resources, a Pkl Template of the CRD should be created.

```bash
pkl eval "package://pkg.pkl-lang.org/pkl-pantry/k8s.contrib.crd@1.0.0#/generate.pkl" -m . -p source="https://raw.githubusercontent.com/argoproj/argo-cd/master/manifests/crds/appproject-crd.yaml"
```
this creates AppProject.pkl in the current Directory.

## Convert Manifest to Pkl
The convert.pkl extends k8s-contribs convert module by allowing non file-based methods of inputs, and by adding the Pkl Template we created earlier as a CRD to the converter.
```bash
pkl eval convert.pkl
```