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
