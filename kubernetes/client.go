package kubernetes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewClient allocates new rest client to interract with k8s
func NewClient(d *schema.ResourceData, terraformVersion string, options client.Options) (client.Client, error) {
	// Config initialization
	cfg, err := GetConfig(d, terraformVersion)
	if err != nil {
		return nil, err
	}
	return client.New(cfg, options)
}
