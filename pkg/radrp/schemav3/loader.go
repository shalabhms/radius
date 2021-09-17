// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package schemav3

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"

	"github.com/xeipuuv/gojsonreference"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// The listing of files below has an ordering to them, because
	// each file may depend on one or more files on the preceding
	// lines.
	//go:embed common-types.json
	//go:embed traits/*.json
	//go:embed traits.json
	//go:embed components/*.json
	//go:embed application.json
	schemaFiles embed.FS

	//go:embed resource-types.json
	manifestFile string

	validators map[string]validator = loadOrPanic()
)

// manifest is the format of the 'resource-types.json' manifest.
type manifest struct {
	Resources map[string]string `json:"resources"`
}

func readManifestOrPanic() manifest {
	manifest := manifest{}
	err := json.Unmarshal([]byte(manifestFile), &manifest)
	if err != nil {
		log.Fatal("Failed to load resource manifest:", err)
	}

	return manifest
}

func GetValidator(resourceType string) (Validator, bool) {
	validator, ok := validators[resourceType]
	if ok {
		return &validator, true
	}

	return nil, false
}

func loadOrPanic() map[string]validator {
	loader := gojsonschema.NewSchemaLoader()
	err := fs.WalkDir(schemaFiles, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(schemaFiles, path)
		if err != nil {
			return fmt.Errorf("cannot read embedded file %s: %w", path, err)
		}
		fileLoader := gojsonschema.NewBytesLoader(data)
		if err = loader.AddSchema( /* url */ "/"+path, fileLoader); err != nil {
			return fmt.Errorf("failed to parse JSON Schema from %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		log.Fatal("Failed to load schemas:", err)
	}

	manifest := readManifestOrPanic()
	validators := map[string]validator{}
	for resourceType, ref := range manifest.Resources {

		// The default logic of the schema loader for references is pretty obtuse. If you give
		// it a reference then it can load from the pool, this is what we want. None of the built-in
		// loaders have this behavior.
		//
		// - Loading from a string will 'poison' the cache because the schema doesn't have a unique reference
		// - Other loaders hit the filesystem/internet which we DO NOT WANT for security reasons.
		workaround := &StrictReferenceLoader{
			Reference: ref,
		}
		schema, err := loader.Compile(workaround)
		if err != nil {
			log.Fatalf("Failed to parse JSON Schema %q: %s", ref, err)
		}
		validators[resourceType] = validator{
			schema:   schema,
			TypeName: resourceType,
		}
	}

	return validators
}

type StrictReferenceLoader struct {
	Reference string
}

func (l *StrictReferenceLoader) JsonSource() interface{} {
	return "/" + l.Reference
}
func (l *StrictReferenceLoader) LoadJSON() (interface{}, error) {
	return nil, errors.New("not supported")
}
func (l *StrictReferenceLoader) JsonReference() (gojsonreference.JsonReference, error) {
	return gojsonreference.NewJsonReference("/" + l.Reference)
}
func (l *StrictReferenceLoader) LoaderFactory() gojsonschema.JSONLoaderFactory {
	return gojsonschema.DefaultJSONLoaderFactory{}
}