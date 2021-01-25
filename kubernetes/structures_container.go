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
	if v, ok := in["optional"]; ok {
		obj.Optional = ptrToBool(v.(bool))
	}
	return obj
}

// ExpandContainerEnv expands container envs
func ExpandContainerEnv(in []interface{}) ([]v1.EnvVar, error) {
	if len(in) == 0 {
		return []v1.EnvVar{}, nil
	}
	envs := make([]v1.EnvVar, len(in))
	for i, c := range in {
		p := c.(map[string]interface{})
		if name, ok := p["name"]; ok {
			envs[i].Name = name.(string)
		}
		if value, ok := p["value"]; ok {
			envs[i].Value = value.(string)
		}
		if v, ok := p["value_from"].([]interface{}); ok && len(v) > 0 {
			var err error
			envs[i].ValueFrom, err = ExpandEnvValueFrom(v)
			if err != nil {
				return envs, err
			}
		}
	}
	return envs, nil
}

// ExpandEnvValueFrom expands env value from diferent sources
func ExpandEnvValueFrom(r []interface{}) (*v1.EnvVarSource, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.EnvVarSource{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.EnvVarSource{}

	var err error
	if v, ok := in["config_map_key_ref"].([]interface{}); ok && len(v) > 0 {
		obj.ConfigMapKeyRef, err = ExpandConfigMapKeyRef(v)
		if err != nil {
			return obj, err
		}
	}
	if v, ok := in["field_ref"].([]interface{}); ok && len(v) > 0 {
		obj.FieldRef, err = ExpandFieldRef(v)
		if err != nil {
			return obj, err
		}
	}
	if v, ok := in["secret_key_ref"].([]interface{}); ok && len(v) > 0 {
		obj.SecretKeyRef = ExpandSecretKeyRef(v)
	}
	if v, ok := in["resource_field_ref"].([]interface{}); ok && len(v) > 0 {
		obj.ResourceFieldRef, err = ExpandResourceFieldRef(v)
		if err != nil {
			return obj, err
		}
	}
	return obj, nil

}

// ExpandConfigMapKeyRef exapnds env value from config map reference
func ExpandConfigMapKeyRef(r []interface{}) (*v1.ConfigMapKeySelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ConfigMapKeySelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ConfigMapKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	if v, ok := in["optional"]; ok {
		obj.Optional = ptrToBool(v.(bool))
	}
	return obj, nil

}

// ExpandFieldRef expands env value from field reference
func ExpandFieldRef(r []interface{}) (*v1.ObjectFieldSelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ObjectFieldSelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ObjectFieldSelector{}

	if v, ok := in["api_version"].(string); ok {
		obj.APIVersion = v
	}
	if v, ok := in["field_path"].(string); ok {
		obj.FieldPath = v
	}
	return obj, nil
}

// ExpandResourceFieldRef expands env value from resource field reference
func ExpandResourceFieldRef(r []interface{}) (*v1.ResourceFieldSelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ResourceFieldSelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ResourceFieldSelector{}

	if v, ok := in["container_name"].(string); ok {
		obj.ContainerName = v
	}
	if v, ok := in["resource"].(string); ok {
		obj.Resource = v
	}
	if v, ok := in["divisor"].(string); ok {
		q, err := resource.ParseQuantity(v)
		if err != nil {
			return obj, err
		}
		obj.Divisor = q
	}
	return obj, nil
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
	if in.Optional != nil {
		att["optional"] = *in.Optional
	}
	return []interface{}{att}
}

// FlattenContainerEnvs translate env vars to terraform structures
func FlattenContainerEnvs(in []v1.EnvVar) []interface{} {
	att := make([]interface{}, len(in))
	for i, v := range in {
		m := map[string]interface{}{}
		if v.Name != "" {
			m["name"] = v.Name
		}
		if v.Value != "" {
			m["value"] = v.Value
		}
		if v.ValueFrom != nil {
			m["value_from"] = FlattenValueFrom(v.ValueFrom)
		}

		att[i] = m
	}
	return att
}

// FlattenValueFrom converts env value from different sources to terraform structures
func FlattenValueFrom(in *v1.EnvVarSource) []interface{} {
	att := make(map[string]interface{})

	if in.ConfigMapKeyRef != nil {
		att["config_map_key_ref"] = FlattenConfigMapKeyRef(in.ConfigMapKeyRef)
	}
	if in.ResourceFieldRef != nil {
		att["resource_field_ref"] = FlattenResourceFieldSelector(in.ResourceFieldRef)
	}
	if in.SecretKeyRef != nil {
		att["secret_key_ref"] = FlattenSecretKeyRef(in.SecretKeyRef)
	}
	if in.FieldRef != nil {
		att["field_ref"] = FlattenObjectFieldSelector(in.FieldRef)
	}
	return []interface{}{att}
}

// FlattenConfigMapKeyRef converts config map key selector into terraform structures
func FlattenConfigMapKeyRef(in *v1.ConfigMapKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	if in.Optional != nil {
		att["optional"] = *in.Optional
	}
	return []interface{}{att}
}

// FlattenResourceFieldSelector converts resource field selector into terraform structures
func FlattenResourceFieldSelector(in *v1.ResourceFieldSelector) []interface{} {
	att := make(map[string]interface{})

	if in.ContainerName != "" {
		att["container_name"] = in.ContainerName
	}
	if in.Resource != "" {
		att["resource"] = in.Resource
	}
	if in.Divisor.String() != "" {
		att["divisor"] = in.Divisor.String()
	}
	return []interface{}{att}
}

// FlattenObjectFieldSelector converts object field selectr into terraform structures
func FlattenObjectFieldSelector(in *v1.ObjectFieldSelector) []interface{} {
	att := make(map[string]interface{})

	if in.APIVersion != "" {
		att["api_version"] = in.APIVersion
	}
	if in.FieldPath != "" {
		att["field_path"] = in.FieldPath
	}
	return []interface{}{att}
}
