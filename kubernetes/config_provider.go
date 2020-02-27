package kubernetes

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mitchellh/go-homedir"
	apimachineryschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// InitConfig initializes k8s configuration
func InitConfig(d *schema.ResourceData) (*rest.Config, error) {
	overrides := &clientcmd.ConfigOverrides{}
	loader := &clientcmd.ClientConfigLoadingRules{}

	if d.Get("load_config_file").(bool) {
		log.Printf("[DEBUG] Trying to load configuration from file")
		if configPath, ok := d.GetOk("config_path"); ok && configPath.(string) != "" {
			path, err := homedir.Expand(configPath.(string))
			if err != nil {
				return nil, err
			}
			log.Printf("[DEBUG] Configuration file is: %s", path)
			loader.ExplicitPath = path

			ctxSuffix := "; default context"

			ctx, ctxOk := d.GetOk("config_context")
			authInfo, authInfoOk := d.GetOk("config_context_auth_info")
			cluster, clusterOk := d.GetOk("config_context_cluster")
			if ctxOk || authInfoOk || clusterOk {
				ctxSuffix = "; overriden context"
				if ctxOk {
					overrides.CurrentContext = ctx.(string)
					ctxSuffix += fmt.Sprintf("; config ctx: %s", overrides.CurrentContext)
					log.Printf("[DEBUG] Using custom current context: %q", overrides.CurrentContext)
				}

				overrides.Context = clientcmdapi.Context{}
				if authInfoOk {
					overrides.Context.AuthInfo = authInfo.(string)
					ctxSuffix += fmt.Sprintf("; auth_info: %s", overrides.Context.AuthInfo)
				}
				if clusterOk {
					overrides.Context.Cluster = cluster.(string)
					ctxSuffix += fmt.Sprintf("; cluster: %s", overrides.Context.Cluster)
				}
				log.Printf("[DEBUG] Using overidden context: %#v", overrides.Context)
			}
		}
	}

	// Overriding with static configuration
	if v, ok := d.GetOk("insecure"); ok {
		overrides.ClusterInfo.InsecureSkipTLSVerify = v.(bool)
	}
	if v, ok := d.GetOk("cluster_ca_certificate"); ok {
		overrides.ClusterInfo.CertificateAuthorityData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := d.GetOk("client_certificate"); ok {
		overrides.AuthInfo.ClientCertificateData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := d.GetOk("host"); ok {
		// Server has to be the complete address of the kubernetes cluster (scheme://hostname:port), not just the hostname,
		// because `overrides` are processed too late to be taken into account by `defaultServerUrlFor()`.
		// This basically replicates what defaultServerUrlFor() does with config but for overrides,
		// see https://github.com/kubernetes/client-go/blob/v12.0.0/rest/url_utils.go#L85-L87
		hasCA := len(overrides.ClusterInfo.CertificateAuthorityData) != 0
		hasCert := len(overrides.AuthInfo.ClientCertificateData) != 0
		defaultTLS := hasCA || hasCert || overrides.ClusterInfo.InsecureSkipTLSVerify
		host, _, err := rest.DefaultServerURL(v.(string), "", apimachineryschema.GroupVersion{}, defaultTLS)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse host: %s", err)
		}

		overrides.ClusterInfo.Server = host.String()
	}
	if v, ok := d.GetOk("username"); ok {
		overrides.AuthInfo.Username = v.(string)
	}
	if v, ok := d.GetOk("password"); ok {
		overrides.AuthInfo.Password = v.(string)
	}
	if v, ok := d.GetOk("client_key"); ok {
		overrides.AuthInfo.ClientKeyData = bytes.NewBufferString(v.(string)).Bytes()
	}
	if v, ok := d.GetOk("token"); ok {
		overrides.AuthInfo.Token = v.(string)
	}

	if v, ok := d.GetOk("exec"); ok {
		exec := &clientcmdapi.ExecConfig{}
		if spec, ok := v.([]interface{})[0].(map[string]interface{}); ok {
			exec.APIVersion = spec["api_version"].(string)
			exec.Command = spec["command"].(string)
			exec.Args = expandStringSlice(spec["args"].([]interface{}))
			for kk, vv := range spec["env"].(map[string]interface{}) {
				exec.Env = append(exec.Env, clientcmdapi.ExecEnvVar{Name: kk, Value: vv.(string)})
			}
		} else {
			return nil, fmt.Errorf("Failed to parse exec")
		}
		overrides.AuthInfo.Exec = exec
	}

	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides)
	cfg, err := cc.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize config: %s", err)
	}

	log.Printf("[INFO] Successfully initialized config")
	return cfg, nil
}

// GetConfig returns REST config for k8s api client
func GetConfig(d *schema.ResourceData, terraformVersion string) (*rest.Config, error) {
	cfg, err := InitConfig(d)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, fmt.Errorf("Failed to initialize config")
	}

	if terraformVersion == "" {
		// Terraform 0.12 introduced this field to the protocol
		// We can therefore assume that if it's missing it's 0.10 or 0.11
		terraformVersion = "0.11+compatible"
	}

	cfg.UserAgent = fmt.Sprintf("HashiCorp/1.0 Terraform/%s", terraformVersion)

	if logging.IsDebugOrHigher() {
		log.Printf("[DEBUG] Enabling HTTP requests/responses tracing")
		cfg.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
			return logging.NewTransport("Kubernetes", rt)
		}
	}

	return cfg, nil
}
