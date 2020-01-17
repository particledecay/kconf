package kubeconfig

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"gopkg.in/oleiade/reflections.v1"
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
func (k *KConf) AddContext(context *clientcmdapi.Context) string {
	// if k.hasContext(context) {

	// }
	// if mainConfig.AddContext(ctx) {
	// 	ctxName = mainConfig.RenameContext(ctxName)
	// 	log.Info().Msgf("Context '%s' already exists, renaming to '%s'", ctxName, newName)
	// }
	return ""
}

// hasContext checks whether the KConf already contains the given Context
func (k *KConf) hasContext(context *clientcmdapi.Context) bool {
	var foundContext *clientcmdapi.Context
	fields := []string{"Cluster", "AuthInfo", "Namespace", "Extensions"}
	for _, ctx := range k.Contexts {
		matched := false
		for _, f := range fields {
			xVal, xErr := reflections.GetField(ctx, f)
			yVal, yErr := reflections.GetField(context, f)
			if xErr == yErr && xVal != yVal {
				matched = false
				break
			} else {
				matched = true
			}
		}
		if matched == true {
			log.Info().Msgf("found the following: %v", ctx)
			foundContext = ctx
			break
		}
	}
	return foundContext == nil
}

func (k *KConf) rename(name string, objType string) string {
	inc := 0
	for {
		if objType == "context" {
			if _, ok := k.Contexts[name]; !ok {
				return name
			} else {
				name = fmt.Sprintf("%s-%d", name, inc)
				inc++
			}
		} else if objType == "cluster" {
			if _, ok := k.Clusters[name]; !ok {
				return name
			} else {
				name = fmt.Sprintf("%s-%d", name, inc)
				inc++
			}
		} else if objType == "user" {
			if _, ok := k.AuthInfos[name]; !ok {
				return name
			} else {
				name = fmt.Sprintf("%s-%d", name, inc)
				inc++
			}
		}
	}
}
