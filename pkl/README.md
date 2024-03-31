# Alternative aproach

Simplify Reader implementation by converting the whole Function Input to Pkl

This would reduce the amount of times Pkl evaluator needs to be invoked. Currently it must be called for each resource/connectionDetails/ready.
Also pkl would likely cache the result for the EvaluatorManager.


## How it's used
```shell
pkl eval convert.pkl
```
input and crds will need to be overwritten by the yaml of the functionInput and a reference to the CRDs used by this composition.


## How to implement it
convert yaml from function input to pkl Template (within the reader implementation)

convert 


the pkl/convert can convert a yaml file to Pkl and provides converters to get it back to yaml