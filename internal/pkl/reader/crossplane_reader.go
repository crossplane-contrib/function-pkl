/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reader

import (
	"context"
	"fmt"
	"net/url"

	"github.com/apple/pkl-go/pkl"
	"sigs.k8s.io/yaml"

	"github.com/crossplane/function-sdk-go/logging"

	"github.com/crossplane-contrib/function-pkl/internal/helper"
)

type CrossplaneReader struct {
	Request      *helper.CompositionRequest
	ReaderScheme string
	Log          logging.Logger
	Ctx          context.Context
}

func (f *CrossplaneReader) Scheme() string {
	return f.ReaderScheme
}

func (f *CrossplaneReader) IsGlobbable() bool {
	return false
}

func (f *CrossplaneReader) HasHierarchicalUris() bool {
	return false
}

// ListElements returns the list of elements at a specified path.
// If HasHierarchicalUris is false, path will be empty and ListElements should return all
// available values.
//
// This method is only called if it is hierarchical and local, or if it is globbable.
func (f *CrossplaneReader) ListElements(_ url.URL) ([]pkl.PathElement, error) {
	out := []pkl.PathElement{
		pkl.NewPathElement("request", false),
	}
	return out, nil
}

var _ pkl.Reader = (*CrossplaneReader)(nil)

type crossplaneModuleReader struct {
	*CrossplaneReader
}

func (f crossplaneModuleReader) IsLocal() bool {
	return true
}

var evaluatorManager pkl.EvaluatorManager = pkl.NewEvaluatorManager()

// TODO find better solution
func Close() error {
	return evaluatorManager.Close()
}

var WithCrossplane = func(crossplaneReader *CrossplaneReader) func(opts *pkl.EvaluatorOptions) {
	return func(opts *pkl.EvaluatorOptions) {
		reader := crossplaneReader
		pkl.WithModuleReader(&crossplaneModuleReader{reader})(opts)
		pkl.WithResourceReader(&crossplaneResourceReader{reader})(opts)
	}
}

func (f *CrossplaneReader) BaseRead(url url.URL) ([]byte, error) {
	switch url.Opaque {
	case "request":
		requestYaml, err := yaml.Marshal(f.Request)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(requestYaml))

		return requestYaml, nil
	default:
		return nil, fmt.Errorf("unsupported path")
	}
}

// Expects an URL like /observed/composition/resource and evaluates the RunFunctionRequest for the state of the desired field and returns it as a pkl file
func (f crossplaneModuleReader) Read(url url.URL) (string, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

var _ pkl.ModuleReader = (*crossplaneModuleReader)(nil)

type crossplaneResourceReader struct {
	*CrossplaneReader
}

func (f crossplaneResourceReader) Read(url url.URL) ([]byte, error) {
	out, err := f.BaseRead(url)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _ pkl.ResourceReader = (*crossplaneResourceReader)(nil)
