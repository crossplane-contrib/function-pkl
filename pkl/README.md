# Alternative aproach

Simplify Reader implementation by converting the whole Function Input to Pkl

This would reduce the amount of times Pkl evaluator needs to be invoked. Currently it must be called for each resource/connectionDetails/ready.
Also pkl would likely cache the result for the EvaluatorManager.

