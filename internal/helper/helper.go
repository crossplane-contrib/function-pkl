package helper

import (
	"fmt"
	"strings"

	"github.com/avarei/function-pkl/input/v1beta1"
)

const coreDefaultPackage string = "package://pkg.pkl-lang.org/github.com/avarei/function-pkl/crossplane@0.0.10"

func ParsePackages(packageList []v1beta1.Package) Packages {
	out := Packages{
		core: coreDefaultPackage,
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
	i := strings.Index(uri, "@")
	if i < 0 {
		return uri
	}
	packageUri, ok := p.packages[uri[:i]]
	if !ok {
		// If no match was found try the full path
		return uri
	}
	return fmt.Sprintf("%s#%s", packageUri, uri[i+1:])
}
