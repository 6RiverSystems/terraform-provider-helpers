package kubernetes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// ProviderFields for configuration
func ProviderFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_HOST", ""),
			Description: "The hostname (in form of URI) of Kubernetes master.",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_USER", ""),
			Description: "The username to use for HTTP basic authentication when accessing the Kubernetes master endpoint.",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_PASSWORD", ""),
			Description: "The password to use for HTTP basic authentication when accessing the Kubernetes master endpoint.",
		},
		"insecure": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_INSECURE", false),
			Description: "Whether server should be accessed without verifying the TLS certificate.",
		},
		"client_certificate": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CLIENT_CERT_DATA", ""),
			Description: "PEM-encoded client certificate for TLS authentication.",
		},
		"client_key": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CLIENT_KEY_DATA", ""),
			Description: "PEM-encoded client certificate key for TLS authentication.",
		},
		"cluster_ca_certificate": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CLUSTER_CA_CERT_DATA", ""),
			Description: "PEM-encoded root certificates bundle for TLS authentication.",
		},
		"config_path": {
			Type:     schema.TypeString,
			Optional: true,
			DefaultFunc: schema.MultiEnvDefaultFunc(
				[]string{
					"KUBE_CONFIG",
					"KUBECONFIG",
				},
				"~/.kube/config"),
			Description: "Path to the kube config file, defaults to ~/.kube/config",
		},
		"config_context": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CTX", ""),
		},
		"config_context_auth_info": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CTX_AUTH_INFO", ""),
			Description: "",
		},
		"config_context_cluster": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_CTX_CLUSTER", ""),
			Description: "",
		},
		"token": {
			Type:        schema.TypeString,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_TOKEN", ""),
			Description: "Token to authenticate an service account",
		},
		"load_config_file": {
			Type:        schema.TypeBool,
			Optional:    true,
			DefaultFunc: schema.EnvDefaultFunc("KUBE_LOAD_CONFIG_FILE", true),
			Description: "Load local kubeconfig.",
		},
		"exec": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_version": {
						Type:     schema.TypeString,
						Required: true,
					},
					"command": {
						Type:     schema.TypeString,
						Required: true,
					},
					"env": {
						Type:     schema.TypeMap,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"args": {
						Type:     schema.TypeList,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
			Description: "",
		},
	}
}
