# TODO List

## Open/WiP
### Prototype

- [x] rename 'full' input to something better
- [x] decide whether to keep composition and resources inputs.
- [x] find a way to load the imports from full Pkl to convert.pkl
      This would remove the need manually specifying them
- [x] [Make Objects neat](https://github.com/apple/pkl-pantry/issues/62)
- [ ] Provide a Pkl Module for Compositions and XRDs, so they can help make Compositions faster.
      Impacted by https://github.com/apple/pkl-pantry/issues/40


- [x] Use PklProject with references instead of static links in code and Pkl Templates.
- [x] Consider making the Dependencies available locally
  - obsolete due to the of the converter. unless there is a nice way to have the package pre-cached.
- [x] Switch from Errors to Results in fn.go

- [x] Cleanup
  - [x] Remove Brainstorming files
  - [x] Improve Variables in pkl/convert.pkl

- [ ] Documentation
  - [x] Update examples
  - [ ] Update README.md
  - [ ] Create Specialized convert.pkl for Migrating XRDs to Pkl

### Road to v1
- [ ] Migrate to github.com/crossplane-contrib
- [ ] Evaluate GitHub Workflows from Template
- [ ] Create and Push Container Image
- [ ] Publish in Upbound MarketPlace

## Completed
### Prototype
- [x] Use one pkl.EvaluatorManager per invocation or globally
- [x] Support Composition Status Updates
- [ ] Disallow Filesystem Access for pkl Modules as that is not needed

### Proof of Concept
- [x] Allow Composition Resources to be loaded from:
    - [x] Uri
    - [x] InlineFile
    ~~- [ ] ConfigMap~~
- [x] Turn Function Input into Pkl Files
    - [x] Implement pkl.Reader Interface.
- [x] Allow Custom Resource Definition Pkl Templates in v1beta1/PklSpec
    - [x] From Uri
    - [ ] ~~From ConfigMap~~
    - [ ] ~~From Inline?~~
        - [x] Use these Pkl Templates to Convert the requested Input to Pkl Files
- [x] Render all Pkl Files
- [x] Allow Ready Status and Connection Details to be Set in the Pkl Manifests
