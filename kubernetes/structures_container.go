package kubernetes

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func expandMapToResourceList(m map[string]interface{}) (*apiv1.ResourceList, error) {
	out := make(apiv1.ResourceList)
	for stringKey, origValue := range m {
		key := apiv1.ResourceName(stringKey)
		var value resource.Quantity

		if v, ok := origValue.(int); ok {
			q := resource.NewQuantity(int64(v), resource.DecimalExponent)
			value = *q
		} else if v, ok := origValue.(string); ok {
			var err error
			value, err = resource.ParseQuantity(v)
			if err != nil {
				return &out, err
			}
		} else {
			return &out, fmt.Errorf("Unexpected value type: %#v", origValue)
		}

		out[key] = value
	}
	return &out, nil
}

// ExpandContainerResourceRequirements converts interface list into resource requirements
func ExpandContainerResourceRequirements(l []interface{}) (*corev1.ResourceRequirements, error) {
	obj := &corev1.ResourceRequirements{}
	if len(l) == 0 || l[0] == nil {
		return obj, nil
	}
	in := l[0].(map[string]interface{})

	fn := func(in []interface{}) (*corev1.ResourceList, error) {
		for _, c := range in {
			p := c.(map[string]interface{})
			if p["cpu"] == "" {
				delete(p, "cpu")
			}
			if p["memory"] == "" {
				delete(p, "memory")
			}
			rl, err := expandMapToResourceList(p)
			if err != nil {
				return rl, err
			}
			return rl, nil
		}
		return nil, nil
	}

	if v, ok := in["limits"].([]interface{}); ok && len(v) > 0 {
		rl, err := fn(v)
		if err != nil {
			return obj, err
		}
		obj.Limits = *rl
	}

	if v, ok := in["requests"].([]interface{}); ok && len(v) > 0 {
		rq, err := fn(v)
		if err != nil {
			return obj, err
		}
		obj.Requests = *rq
	}

	return obj, nil
}

// ExpandSecretKeyRef converts terraform list into secret key selector
func ExpandSecretKeyRef(r []interface{}) *v1.SecretKeySelector {
	if len(r) == 0 || r[0] == nil {
		return &v1.SecretKeySelector{}
	}
	in := r[0].(map[string]interface{})
	obj := &v1.SecretKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	return obj
}

// FlattenResourceList converts resource list into terraform list
func FlattenResourceList(l apiv1.ResourceList) map[string]string {
	m := make(map[string]string)
	for k, v := range l {
		m[string(k)] = v.String()
	}
	return m
}

// FlattenContainerResourceRequirements converts resource requirements into interface list
func FlattenContainerResourceRequirements(in corev1.ResourceRequirements) ([]interface{}, error) {
	att := make(map[string]interface{})
	if len(in.Limits) > 0 {
		att["limits"] = []interface{}{FlattenResourceList(in.Limits)}
	}
	if len(in.Requests) > 0 {
		att["requests"] = []interface{}{FlattenResourceList(in.Requests)}
	}
	return []interface{}{att}, nil
}

// FlattenSecretKeyRef converts selector into terraform list
func FlattenSecretKeyRef(in *v1.SecretKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	return []interface{}{att}
}
