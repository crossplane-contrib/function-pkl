# TODO List

- [ ] Use one pkl.EvaluatorManager per invocation or globally
- [ ] Disallow Filesystem Access for pkl Modules as that is not needed
- [ ] Allow Composition Resources to be loaded from:
    - [x] Uri
    - [x] InlineFile
    - [ ] ConfigMap
- [ ] Turn Function Input into Pkl Files
    - [x] Implement pkl.Reader Interface.
- [ ] Allow Custom Resource Definition Pkl Templates in v1beta1/PklSpec
    - [x] From Uri
    - [ ] ~~From ConfigMap~~
    - [ ] ~~From Inline?~~
        - [x] Use these Pkl Templates to Convert the requested Input to Pkl Files
- [x] Render all Pkl Files
- [ ] Use PklProject with references instead of static links in code and Pkl Templates.
- [ ] Allow Ready Status and Connection Details to be Set in the Pkl Manifests
- [ ] Attempt to Remove ApiVersion and Kind fields from Function input
- [ ] Cleanup
  - [ ] Remove Brainstorming files
  - [ ] Improve Variables in pkl/convert.pkl

- [ ] Documentation
  - [ ] Update examples
  - [ ] Update README.md
