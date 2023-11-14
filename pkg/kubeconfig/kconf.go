package kubeconfig

import (
	"fmt"
	"os"
	"path"
	"reflect"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Configurable can validate and modify its own kubeconfig attributes
type Configurable interface {
	AddContext() bool
	AddCluster() bool
	AddUser() bool

	hasContext() bool
	hasCluster() bool
	hasUser() bool
	rename() string
}

// KConf is a kubeconfig mixed with some useful functions
type KConf struct {
	clientcmdapi.Config
}

// MainConfigPath is the file path to the main config
var MainConfigPath string

func init() {
	MainConfigPath = path.Join(os.Getenv("HOME"), ".kube", "config")
}

// AddContext attempts to add a Context and return the resulting name
func (k *KConf) AddContext(name string, context *clientcmdapi.Context) string {
	// context exists and is identical
	if _, ok := k.Contexts[name]; ok && k.hasContext(context) {
		return ""
	}
	name, _ = k.rename(name, "context")

	k.Config.Contexts[name] = context
	return name
}

// IsEqualContext checks whether two Contexts are equal
func IsEqualContext(ctx1, ctx2 *clientcmdapi.Context) bool {
	if ctx1.Cluster != ctx2.Cluster ||
		ctx1.AuthInfo != ctx2.AuthInfo ||
		!reflect.DeepEqual(ctx1.Extensions, ctx2.Extensions) {
		return false
	}
	return true
}

// hasContext checks whether the KConf already contains the given Context
func (k *KConf) hasContext(context *clientcmdapi.Context) bool {
	for _, ctx := range k.Contexts {
		if IsEqualContext(ctx, context) {
			return true
		}
	}
	return false
}

// AddCluster attempts to add a Cluster and return the resulting name
func (k *KConf) AddCluster(name string, cluster *clientcmdapi.Cluster) string {
	if k.hasCluster(cluster) {
		return ""
	}
	name, _ = k.rename(name, "cluster")

	k.Config.Clusters[name] = cluster
	return name
}

// IsEqualCluster checks whether two Clusters are equal
func IsEqualCluster(cls1, cls2 *clientcmdapi.Cluster) bool {
	if cls1.CertificateAuthority != cls2.CertificateAuthority ||
		string(cls1.CertificateAuthorityData) != string(cls2.CertificateAuthorityData) ||
		cls1.Server != cls2.Server ||
		cls1.InsecureSkipTLSVerify != cls2.InsecureSkipTLSVerify ||
		!reflect.DeepEqual(cls1.Extensions, cls2.Extensions) {
		return false
	}
	return true
}

// hasCluster checks whether the KConf already contains the given Cluster
func (k *KConf) hasCluster(cluster *clientcmdapi.Cluster) bool {
	for _, cls := range k.Clusters {
		if IsEqualCluster(cls, cluster) {
			return true
		}
	}
	return false
}

// AddUser attempts to add an AuthInfo and return the resulting name
func (k *KConf) AddUser(name string, user *clientcmdapi.AuthInfo) string {
	if k.hasUser(user) {
		return ""
	}
	name, _ = k.rename(name, "user")

	k.Config.AuthInfos[name] = user
	return name
}

// IsEqualUser checks whether two Users are equal
func IsEqualUser(user1, user2 *clientcmdapi.AuthInfo) bool {
	if user1.ClientCertificate != user2.ClientCertificate ||
		string(user1.ClientCertificateData) != string(user2.ClientCertificateData) ||
		user1.ClientKey != user2.ClientKey ||
		string(user1.ClientKeyData) != string(user2.ClientKeyData) ||
		user1.Token != user2.Token ||
		user1.TokenFile != user2.TokenFile ||
		user1.Impersonate != user2.Impersonate ||
		user1.ImpersonateUID != user2.ImpersonateUID ||
		!reflect.DeepEqual(user1.ImpersonateGroups, user2.ImpersonateGroups) ||
		user1.Username != user2.Username ||
		user1.Password != user2.Password ||
		!reflect.DeepEqual(user1.AuthProvider, user2.AuthProvider) ||
		!reflect.DeepEqual(user1.Exec, user2.Exec) {
		return false
	}
	return true
}

// hasUser checks whether the KConf already contains the given AuthInfo
func (k *KConf) hasUser(user *clientcmdapi.AuthInfo) bool {
	for _, u := range k.AuthInfos {
		if IsEqualUser(u, user) {
			return true
		}
	}
	return false
}

func (k *KConf) rename(name string, objType string) (string, error) {
	inc := 1
	origName := name
	for {
		switch objType {
		case "context":
			if _, ok := k.Contexts[name]; !ok {
				return name, nil
			}
		case "cluster":
			if _, ok := k.Clusters[name]; !ok {
				return name, nil
			}
		case "user":
			if _, ok := k.AuthInfos[name]; !ok {
				return name, nil
			}
		default:
			return "", fmt.Errorf("unrecognized type '%s'", objType)
		}

		name = fmt.Sprintf("%s-%d", origName, inc)
		inc++
	}
}

// MoveContext renames an existing context
func (k *KConf) MoveContext(oldName, newName string) error {
	ctx, ok := k.Contexts[oldName]
	if !ok {
		return fmt.Errorf("no existing context named '%s'", oldName)
	}
	if _, ok := k.Contexts[newName]; ok {
		return fmt.Errorf("context named '%s' already exists", newName)
	}

	k.Contexts[newName] = ctx
	delete(k.Contexts, oldName)

	return nil
}

// contextsWithUser returns a slice of context names that use a specific AuthInfo
func (k *KConf) contextsWithUser(name string) []string {
	contexts := []string{}

	for ctxName, ctx := range k.Contexts {
		if ctx.AuthInfo == name {
			contexts = append(contexts, ctxName)
		}
	}

	return contexts
}
