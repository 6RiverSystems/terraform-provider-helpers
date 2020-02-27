package kubernetes

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// LocalObjectReferenceSchema produces schema for local reference object
func LocalObjectReferenceSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Description: "Contains enough information to let you locate the referenced object inside the same namespace",
		Required:    required,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Description:  "Name of the referent",
					Required:     true,
					ForceNew:     true,
					Computed:     false,
					ValidateFunc: ValidateName,
				},
			},
		},
	}
}
