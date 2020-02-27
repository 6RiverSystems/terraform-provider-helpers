package kubernetes

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExpandLocalObjectReference is map[string]interface{} to reference transformation
func ExpandLocalObjectReference(in []interface{}) corev1.LocalObjectReference {
	ref := corev1.LocalObjectReference{}
	if len(in) < 1 {
		return ref
	}

	m := in[0].(map[string]interface{})
	if v, ok := m["name"]; ok {
		ref.Name = v.(string)
	}

	return ref
}

// FlattenLocalObjectReference is local object reference to map transformation
func FlattenLocalObjectReference(ref corev1.LocalObjectReference) []interface{} {
	m := make(map[string]interface{})
	m["name"] = ref.Name
	return []interface{}{m}
}

// IDParts parses ID returning namespace, name, and optional error
func IDParts(id string) (string, string, error) {
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		err := fmt.Errorf("unexpected ID format (%q) expected %q", id, "namespace/name")
		return "", "", err
	}

	return parts[0], parts[1], nil
}

// BuildID composes ID from namespace and resource name
func BuildID(meta metav1.ObjectMeta) string {
	return meta.Namespace + "/" + meta.Name
}

func expandStringSlice(s []interface{}) []string {
	result := make([]string, len(s), len(s))
	for k, v := range s {
		// Handle the Terraform parser bug which turns empty strings in lists to nil.
		if v == nil {
			result[k] = ""
		} else {
			result[k] = v.(string)
		}
	}
	return result
}
