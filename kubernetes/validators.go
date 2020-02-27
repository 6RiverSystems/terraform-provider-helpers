package kubernetes

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
	apiValidation "k8s.io/apimachinery/pkg/api/validation"
	utilValidation "k8s.io/apimachinery/pkg/util/validation"
)

// ValidateAnnotations validatesresource annotation
func ValidateAnnotations(value interface{}, key string) (ws []string, es []error) {
	m := value.(map[string]interface{})
	for k := range m {
		errors := utilValidation.IsQualifiedName(strings.ToLower(k))
		if len(errors) > 0 {
			for _, e := range errors {
				es = append(es, fmt.Errorf("%s (%q) %s", key, k, e))
			}
		}
	}
	return
}

// ValidateName validages the string is valid name
func ValidateName(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)
	errors := apiValidation.NameIsDNSSubdomain(v, false)
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

// ValidateLabels validates resource labels
func ValidateLabels(value interface{}, key string) (ws []string, es []error) {
	m := value.(map[string]interface{})
	for k, v := range m {
		for _, msg := range utilValidation.IsQualifiedName(k) {
			es = append(es, fmt.Errorf("%s (%q) %s", key, k, msg))
		}
		val, isString := v.(string)
		if !isString {
			es = append(es, fmt.Errorf("%s.%s (%#v): Expected value to be string", key, k, v))
			return
		}
		for _, msg := range utilValidation.IsValidLabelValue(val) {
			es = append(es, fmt.Errorf("%s (%q) %s", key, val, msg))
		}
	}
	return
}

// ValidateGenerateName validate resource generate name
func ValidateGenerateName(value interface{}, key string) (ws []string, es []error) {
	v := value.(string)

	errors := apiValidation.NameIsDNSLabel(v, true)
	if len(errors) > 0 {
		for _, err := range errors {
			es = append(es, fmt.Errorf("%s %s", key, err))
		}
	}
	return
}

// ValidateResourceQuantity validate resource quantity
func ValidateResourceQuantity(value interface{}, key string) (ws []string, es []error) {
	if v, ok := value.(string); ok {
		_, err := resource.ParseQuantity(v)
		if err != nil {
			es = append(es, fmt.Errorf("%s.%s : %s", key, v, err))
		}
	}
	return
}
