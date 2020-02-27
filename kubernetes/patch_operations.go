package kubernetes

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

func diffStringMap(pathPrefix string, oldV, newV map[string]interface{}) PatchOperations {
	ops := make([]PatchOperation, 0, 0)

	pathPrefix = strings.TrimRight(pathPrefix, "/")

	// If old value was empty, just create the object
	if len(oldV) == 0 {
		ops = append(ops, &AddOperation{
			Path:  pathPrefix,
			Value: newV,
		})
		return ops
	}

	// This is suboptimal for adding whole new map from scratch
	// or deleting the whole map, but it's actually intention.
	// There may be some other map items managed outside of TF
	// and we don't want to touch these.

	for k := range oldV {
		if _, ok := newV[k]; ok {
			continue
		}
		ops = append(ops, &RemoveOperation{
			Path: pathPrefix + "/" + escapeJSONPointer(k),
		})
	}

	for k, v := range newV {
		newValue := v.(string)

		if oldValue, ok := oldV[k].(string); ok {
			if oldValue == newValue {
				continue
			}

			ops = append(ops, &ReplaceOperation{
				Path:  pathPrefix + "/" + escapeJSONPointer(k),
				Value: newValue,
			})
			continue
		}

		ops = append(ops, &AddOperation{
			Path:  pathPrefix + "/" + escapeJSONPointer(k),
			Value: newValue,
		})
	}

	return ops
}

// escapeJSONPointer escapes string per RFC 6901
// so it can be used as path in JSON patch operations
func escapeJSONPointer(path string) string {
	path = strings.Replace(path, "~", "~0", -1)
	path = strings.Replace(path, "/", "~1", -1)
	return path
}

// PatchOperations is array of patch operations
type PatchOperations []PatchOperation

// MarshalJSON marshals operations to json
func (po PatchOperations) MarshalJSON() ([]byte, error) {
	var v []PatchOperation = po
	return json.Marshal(v)
}

// Equal compares operations
func (po PatchOperations) Equal(ops []PatchOperation) bool {
	var v []PatchOperation = po

	sort.Slice(v, sortByPathAsc(v))
	sort.Slice(ops, sortByPathAsc(ops))

	return reflect.DeepEqual(v, ops)
}

func sortByPathAsc(ops []PatchOperation) func(i, j int) bool {
	return func(i, j int) bool {
		return ops[i].GetPath() < ops[j].GetPath()
	}
}

// PatchOperation interface for Add, Replace, Remove operations
type PatchOperation interface {
	MarshalJSON() ([]byte, error)
	// GetPath erer
	GetPath() string
}

// ReplaceOperation to replace resource item
type ReplaceOperation struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
	Op    string      `json:"op"`
}

// GetPath returns patch path
func (o *ReplaceOperation) GetPath() string {
	return o.Path
}

// MarshalJSON serializes struct to JSON
func (o *ReplaceOperation) MarshalJSON() ([]byte, error) {
	o.Op = "replace"
	return json.Marshal(*o)
}

func (o *ReplaceOperation) String() string {
	b, _ := o.MarshalJSON()
	return string(b)
}

// AddOperation to resource item
type AddOperation struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
	Op    string      `json:"op"`
}

// GetPath returns patch path
func (o *AddOperation) GetPath() string {
	return o.Path
}

// MarshalJSON serializes struct to JSON
func (o *AddOperation) MarshalJSON() ([]byte, error) {
	o.Op = "add"
	return json.Marshal(*o)
}

func (o *AddOperation) String() string {
	b, _ := o.MarshalJSON()
	return string(b)
}

// RemoveOperation removes items from resource
type RemoveOperation struct {
	Path string `json:"path"`
	Op   string `json:"op"`
}

// GetPath returns patch path
func (o *RemoveOperation) GetPath() string {
	return o.Path
}

// MarshalJSON serializes struct to JSON
func (o *RemoveOperation) MarshalJSON() ([]byte, error) {
	o.Op = "remove"
	return json.Marshal(*o)
}

func (o *RemoveOperation) String() string {
	b, _ := o.MarshalJSON()
	return string(b)
}
