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
	if k.hasContext(context) {
		return ""
	}
	name, _ = k.rename(name, "context")

	k.Config.Contexts[name] = context
	return name
}

// hasContext checks whether the KConf already contains the given Context
func (k *KConf) hasContext(context *clientcmdapi.Context) bool {
	var foundContext bool
	for _, ctx := range k.Contexts {
		if reflect.DeepEqual(ctx, context) == true {
			foundContext = true
			break
		}
	}
	return foundContext
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

// hasCluster checks whether the KConf already contains the given Cluster
func (k *KConf) hasCluster(cluster *clientcmdapi.Cluster) bool {
	var foundCluster bool
	for _, cls := range k.Clusters {
		if reflect.DeepEqual(cls, cluster) == true {
			foundCluster = true
			break
		}
	}
	return foundCluster
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

// hasUser checks whether the KConf already contains the given AuthInfo
func (k *KConf) hasUser(user *clientcmdapi.AuthInfo) bool {
	var foundUser bool
	for _, u := range k.AuthInfos {
		if reflect.DeepEqual(u, user) == true {
			foundUser = true
			break
		}
	}
	return foundUser
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
