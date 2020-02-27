package kubernetes

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExpandMetadata converts terraform object to k8s ObjectMeta
func ExpandMetadata(in []interface{}) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{}
	if len(in) < 1 {
		return meta
	}
	m := in[0].(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = expandStringMap(m["annotations"].(map[string]interface{}))
	}

	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = expandStringMap(m["labels"].(map[string]interface{}))
	}

	if v, ok := m["generate_name"]; ok {
		meta.GenerateName = v.(string)
	}
	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}
	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	return meta
}

// FlattenMetadata falttens k8s metadata into map[string]interface{}
func FlattenMetadata(meta metav1.ObjectMeta, d *schema.ResourceData, metaPrefix ...string) []interface{} {
	m := make(map[string]interface{})
	prefix := ""
	if len(metaPrefix) > 0 {
		prefix = metaPrefix[0]
	}
	configAnnotations := d.Get(prefix + "metadata.0.annotations").(map[string]interface{})
	m["annotations"] = removeInternalKeys(meta.Annotations, configAnnotations)
	if meta.GenerateName != "" {
		m["generate_name"] = meta.GenerateName
	}
	configLabels := d.Get(prefix + "metadata.0.labels").(map[string]interface{})
	m["labels"] = removeInternalKeys(meta.Labels, configLabels)
	m["name"] = meta.Name
	m["resource_version"] = meta.ResourceVersion
	m["self_link"] = meta.SelfLink
	m["uid"] = fmt.Sprintf("%v", meta.UID)
	m["generation"] = meta.Generation

	if meta.Namespace != "" {
		m["namespace"] = meta.Namespace
	}

	return []interface{}{m}
}

// PatchMetadata produces set of patch operations for metadata
func PatchMetadata(keyPrefix, pathPrefix string, d *schema.ResourceData) PatchOperations {
	ops := make([]PatchOperation, 0, 0)
	if d.HasChange(keyPrefix + "annotations") {
		oldV, newV := d.GetChange(keyPrefix + "annotations")
		diffOps := diffStringMap(pathPrefix+"annotations", oldV.(map[string]interface{}), newV.(map[string]interface{}))
		ops = append(ops, diffOps...)
	}
	if d.HasChange(keyPrefix + "labels") {
		oldV, newV := d.GetChange(keyPrefix + "labels")
		diffOps := diffStringMap(pathPrefix+"labels", oldV.(map[string]interface{}), newV.(map[string]interface{}))
		ops = append(ops, diffOps...)
	}
	return ops
}

func removeInternalKeys(m map[string]string, d map[string]interface{}) map[string]string {
	for k := range m {
		if isInternalKey(k) && !isKeyInMap(k, d) {
			delete(m, k)
		}
	}
	return m
}

func isKeyInMap(key string, d map[string]interface{}) bool {
	if d == nil {
		return false
	}
	for k := range d {
		if k == key {
			return true
		}
	}
	return false
}

func isInternalKey(annotationKey string) bool {
	u, err := url.Parse("//" + annotationKey)
	if err == nil && strings.HasSuffix(u.Hostname(), "kubernetes.io") {
		return true
	}

	// Specific to DaemonSet annotations, generated & controlled by the server.
	if strings.Contains(annotationKey, "deprecated.daemonset.template.generation") {
		return true
	}

	return false
}

func expandStringMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}
