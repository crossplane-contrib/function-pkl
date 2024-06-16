# TODO List

## Open/WiP
### Prototype

- [ ] Create a Pkl Module for creating Compositions and XRDs, to help bootstrapping them.
      Impacted, but not blocked by https://github.com/apple/pkl-pantry/issues/40 since Composition uses 

- [ ] Disallow Filesystem Access for pkl Modules as that is not needed

- [ ] Documentation
  - [x] Update examples
  - [x] Update README.md
  - [ ] Create Specialized convert.pkl for Migrating XRDs to Pkl

### Road to v1
- [ ] Migrate to github.com/crossplane-contrib
- [ ] Evaluate GitHub Workflows from Template
- [x] Create and Push Container Image
- [ ] Publish in Upbound MarketPlace
- [ ] Remove this TODO list and create Issues for the remaining points.

## Completed
### Prototype
- [x] Use one pkl.EvaluatorManager per invocation or globally
- [x] Support Composition Status Updates
- [x] Cleanup
  - [x] Remove Brainstorming files
  - [x] Improve Variables in pkl/convert.pkl
- [x] rename 'full' input to something better
- [x] decide whether to keep composition and resources inputs.
- [x] find a way to load the imports from full Pkl to convert.pkl
      This would remove the need manually specifying them
- [x] [Make Objects neat](https://github.com/apple/pkl-pantry/issues/62)
- [x] Use PklProject with references instead of static links in code and Pkl Templates.
- [x] Consider making the Dependencies available locally
  - obsolete due to the of the converter. unless there is a nice way to have the package pre-cached.
- [x] Switch from Errors to Results in fn.go

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
