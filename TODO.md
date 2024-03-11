# TODO List

- [ ] Use pkl.EvaluatorManager
- [ ] Disallow Filesystem Access for pkl Modules as that is not needed
- [ ] Allow Composition Resources to be loaded from:
    - [x] Uri
    - [x] InlineFile
    - [ ] ConfigMap
- [ ] Turn Function Input into Pkl Files
    - [ ] Implement pkl.Reader Interface.  
          This allows import() and read() statements to accept `crossplane:/observed/composite/resource`.  
        - [ ] Create connectionDetails Pkl Template
        - [ ] Create Ready Pkl Template
        - [ ] Allow Custom Resource Definition Pkl Templates in v1beta1/PklSpec
            - [x] From Uri
            - [ ] From ConfigMap
            - [x] From Inline?
        - [ ] Use these Pkl Templates to Convert the requested Input to Pkl Files
            - [ ] Create k8s-contrib/Convert.pkl equivalent for connectionDetails
            - [ ] Create k8s-contrib/Convert.pkl equivalent for Ready
            - [x] Use k8s-contrib/Convert.pkl
- [x] Render all Pkl Files
