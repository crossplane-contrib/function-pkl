package main

import (
	"io/fs"
	"net/url"
	"strings"

	"github.com/apple/pkl-go/pkl"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

type crossplaneReader struct {
	request *fnv1beta1.RunFunctionRequest
	scheme  string
}

func (f *crossplaneReader) Scheme() string {
	return f.scheme
}

func (f *crossplaneReader) IsGlobbable() bool {
	return true
}

// e.g. crossplane:/observed/composition/resource
func (f *crossplaneReader) HasHierarchicalUris() bool {
	return true
}

// e.g. crossplane:/observed/composition/
// TODO
func (f *crossplaneReader) ListElements(url url.URL) ([]pkl.PathElement, error) {
	path := strings.TrimSuffix(strings.TrimPrefix(url.Path, "/"), "/")

	entries, err := fs.ReadDir(f.fs, path)
	if err != nil {
		return nil, err
	}
	var ret []pkl.PathElement
	for _, entry := range entries {
		// copy Pkl's built-in `file` ModuleKey and don't follow symlinks.
		if entry.Type()&fs.ModeSymlink != 0 {
			continue
		}
		ret = append(ret, pkl.NewPathElement(entry.Name(), entry.IsDir()))
	}
	return ret, nil
}

var _ pkl.Reader = (*crossplaneReader)(nil)

type crossplaneModuleReader struct {
	*crossplaneReader
}

func (f crossplaneModuleReader) IsLocal() bool {
	return true
}

var WithCrossplane = func(req *fnv1beta1.RunFunctionRequest, scheme string) func(opts *pkl.EvaluatorOptions) {
	return func(opts *pkl.EvaluatorOptions) {
		reader := &crossplaneReader{request: req, scheme: scheme}
		pkl.WithModuleReader(&crossplaneModuleReader{reader})(opts)
		pkl.WithResourceReader(&crossplaneResourceReader{reader})(opts)
	}
}

func (f crossplaneModuleReader) Read(url url.URL) (string, error) {
	path := strings.TrimSuffix(strings.TrimPrefix(url.Path, "/"), "/")
	pathElements := strings.Split(path, "/")

	var state *fnv1beta1.State
	switch pathElements[0] {
	case "observed":
		state = f.request.GetObserved()
	case "desired":
		state = f.request.GetDesired()
	default:
		// ERR
	}

	pathElements = pathElements[:1]
	var resource *fnv1beta1.Resource
	var isComposition = false
	switch pathElements[0] {
	case "composition":
		isComposition = true
		resource = state.GetComposite()
	case "resources":
		resource = state.GetResources()[pathElements[1]]
		pathElements = pathElements[:1]
	default:
		// ERR
	}

	pathElements = pathElements[:1]
	switch pathElements[0] {
	case "resource":
		resource.GetResource() // convert using
	case "connectionDetails":
		resource.GetConnectionDetails()
	case "ready":
		resource.GetReady()
	}

	contents, err := fs.ReadFile(f.fs, strings.TrimPrefix(url.Path, "/"))
	return string(contents), err
}

var _ pkl.ModuleReader = (*crossplaneModuleReader)(nil)

type crossplaneResourceReader struct {
	*crossplaneReader
}

func (f crossplaneResourceReader) Read(url url.URL) ([]byte, error) {
	return fs.ReadFile(f.fs, strings.TrimPrefix(url.Path, "/"))
}

var _ pkl.ResourceReader = (*crossplaneResourceReader)(nil)
