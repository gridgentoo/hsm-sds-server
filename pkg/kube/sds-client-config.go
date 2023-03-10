package kube

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

const (
	contextName  = "context0"
	clusterName  = "cluster0"
	authInfoName = "authInfo0"
)

var _ clientcmd.ClientConfig = &clientConfig{}

// clientConfig is a utility that allows construction of a k8s ClientConfig from
// a k8s rest.Config
type clientConfig struct {
	restConfig rest.Config
}

// NewClientConfigForRestConfig creates a new k8s clientcmd.ClientConfig from the given rest.Config.
func NewClientConfigForRestConfig(restConfig *rest.Config) clientcmd.ClientConfig {
	return &clientConfig{
		restConfig: *restConfig,
	}
}

func (c *clientConfig) RawConfig() (api.Config, error) {
	cfg := api.Config{
		Kind:        "Config",
		APIVersion:  "v1",
		Preferences: api.Preferences{},
		Clusters: map[string]*api.Cluster{
			clusterName: newCluster(&c.restConfig),
		},
		AuthInfos: map[string]*api.AuthInfo{
			authInfoName: newAuthInfo(&c.restConfig),
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: authInfoName,
			},
		},
		CurrentContext: contextName,
	}

	return cfg, nil
}

func (c *clientConfig) ClientConfig() (*rest.Config, error) {
	return c.copyRestConfig(), nil
}

func (c *clientConfig) Namespace() (string, bool, error) {
	return "default", false, nil
}

func (c *clientConfig) ConfigAccess() clientcmd.ConfigAccess {
	return nil
}

func (c *clientConfig) copyRestConfig() *rest.Config {
	out := c.restConfig
	return &out
}

func newAuthInfo(restConfig *rest.Config) *api.AuthInfo {
	return &api.AuthInfo{
		ClientCertificate:     restConfig.CertFile,
		ClientCertificateData: restConfig.CertData,
		ClientKey:             restConfig.KeyFile,
		ClientKeyData:         restConfig.KeyData,
		Token:                 restConfig.BearerToken,
		TokenFile:             restConfig.BearerTokenFile,
		Impersonate:           restConfig.Impersonate.UserName,
		ImpersonateGroups:     restConfig.Impersonate.Groups,
		ImpersonateUserExtra:  restConfig.Impersonate.Extra,
		Username:              restConfig.Username,
		Password:              restConfig.Password,
		AuthProvider:          restConfig.AuthProvider,
		Exec:                  restConfig.ExecProvider,
	}
}

func newCluster(restConfig *rest.Config) *api.Cluster {
	return &api.Cluster{
		Server:                   restConfig.Host,
		TLSServerName:            restConfig.ServerName,
		InsecureSkipTLSVerify:    restConfig.Insecure,
		CertificateAuthority:     restConfig.CAFile,
		CertificateAuthorityData: restConfig.CAData,
	}
}
