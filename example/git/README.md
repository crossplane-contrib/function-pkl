# Using Raw Git References
Instead of releasing a Pkl Package this approach reads directly from a git repository.

Currently this method does ot support PklProject dependencies, meaning all dependencies must be declared explicitely.

```pkl
// e.g.
amends "package://pkg.pkl-lang.org/github.com/crossplane-contrib/function-pkl/crossplane.contrib@0.0.1#/CompositionResponse.pkl"

// instead of
amends "@crossplane.contrib/CompositionResponse.pkl"
```

For more information about Generating the Pkl Modules, XRDs, and Compositions see [DEVELOP.md](../../pkl/crossplane.contrib.example/DEVELOP.md)
