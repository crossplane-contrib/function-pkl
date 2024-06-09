package helper

import (
	"fmt"
	"strings"

	"github.com/avarei/function-pkl/input/v1beta1"
)

const coreDefaultPackage string = "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.16"

func ParsePackages(packageList []v1beta1.Package) *Packages {
	out := &Packages{
		core:     coreDefaultPackage,
		packages: make(map[string]string),
	}
	for _, p := range packageList {
		out.packages[p.Name] = p.Uri
		if p.Core {
			out.core = p.Uri
		}
	}
	return out
}

type Packages struct {
	packages map[string]string
	core     string
}

func (p Packages) GetCoreUri() string {
	return p.core
}

func (p Packages) ParseCoreUri(uri string) string {
	return fmt.Sprintf("%s#%s", p.GetCoreUri(), uri)
}

func (p Packages) ParseUri(uri string) string {
	if !strings.HasPrefix(uri, "@") {
		return uri
	}
	i := strings.Index(uri, "/")
	if i < 0 {
		// NOTE: this maybe should even error?
		return uri
	}
	packageUri, ok := p.packages[uri[1:i]]
	if !ok {
		// If no match was found try the full path
		return uri
	}
	return fmt.Sprintf("%s#%s", packageUri, uri[i:])
}
