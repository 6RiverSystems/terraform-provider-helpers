package kubernetes

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// ResourcesField generates schema fileds for container resources field
func ResourcesField() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"limits": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: "Describes the maximum amount of compute resources allowed. More info: http://kubernetes.io/docs/user-guide/compute-resources/",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     ValidateResourceQuantity,
						DiffSuppressFunc: suppressEquivalentResourceQuantity,
					},
					"memory": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     ValidateResourceQuantity,
						DiffSuppressFunc: suppressEquivalentResourceQuantity,
					},
				},
			},
		},
		"requests": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cpu": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     ValidateResourceQuantity,
						DiffSuppressFunc: suppressEquivalentResourceQuantity,
					},
					"memory": {
						Type:             schema.TypeString,
						Optional:         true,
						Computed:         true,
						ValidateFunc:     ValidateResourceQuantity,
						DiffSuppressFunc: suppressEquivalentResourceQuantity,
					},
				},
			},
		},
	}
}

// ContainerEnvFields provides schema definition for container Env array
func ContainerEnvFields(isUpdatable bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    !isUpdatable,
			Description: "Name of the environment variable. Must be a C_IDENTIFIER",
		},
		"value": {
			Type:        schema.TypeString,
			ForceNew:    !isUpdatable,
			Optional:    true,
			Description: `Variable references $(VAR_NAME) are expanded using the previous defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. The $(VAR_NAME) syntax can be escaped with a double $$, ie: $$(VAR_NAME). Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "".`,
		},
		"value_from": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Source for the environment variable's value",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"config_map_key_ref": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Selects a key of a ConfigMap.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "The key to select.",
								},
								"name": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
								},
								"optional": {
									Type:        schema.TypeBool,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "Specify whether the ConfigMap or its key must be defined.",
								},
							},
						},
					},
					"field_ref": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Selects a field of the pod: supports metadata.name, metadata.namespace, metadata.labels, metadata.annotations, spec.nodeName, spec.serviceAccountName, status.podIP.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"api_version": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Default:     "v1",
									Description: `Version of the schema the FieldPath is written in terms of, defaults to "v1".`,
								},
								"field_path": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "Path of the field to select in the specified API version",
								},
							},
						},
					},
					"resource_field_ref": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Selects a resource of the container: only resources limits and requests (limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"container_name": {
									Type:     schema.TypeString,
									Optional: true,
									ForceNew: !isUpdatable,
								},
								"divisor": {
									Type:             schema.TypeString,
									Optional:         true,
									Default:          "1",
									ValidateFunc:     ValidateResourceQuantity,
									DiffSuppressFunc: suppressEquivalentResourceQuantity,
								},
								"resource": {
									Type:        schema.TypeString,
									Required:    true,
									ForceNew:    !isUpdatable,
									Description: "Resource to select",
								},
							},
						},
					},
					"secret_key_ref": {
						Type:        schema.TypeList,
						Optional:    true,
						MaxItems:    1,
						Description: "Selects a key of a secret in the pod's namespace.",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "The key of the secret to select from. Must be a valid secret key.",
								},
								"name": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
								},
								"optional": {
									Type:        schema.TypeBool,
									Optional:    true,
									ForceNew:    !isUpdatable,
									Description: "Specify whether the Secret or its key must be defined.",
								},
							},
						},
					},
				},
			},
		},
	}
}
