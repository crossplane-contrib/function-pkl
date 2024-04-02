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

This composition-function implements the [Reader interface](https://github.com/apple/pkl-go/blob/main/pkl/reader.go)

The following methods are currently implemented:

### State
gives a CompositionInput in pkl format. To be used in functions.
```pkl
import "../CompositionInput.pkl"
local state = import("crossplane:state") as CompositionInput
```

### Input
gives the CompositionInput as a Yaml String
```pkl
read("crossplane:input")
```
### CRDs
returns all required CRDs as Pkl Templates.
```pkl
import*("crossplane:crds")
```
