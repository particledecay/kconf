package kubeconfig

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Save writes out the updated main config to the filesystem
func (k *KConf) Save() error {
	err := clientcmd.WriteToFile(k.Config, MainConfigPath)
	if err != nil {
		log.Debug().
			Err(err).
			Msgf("error while writing main config: %v", err)
		return err
	}
	return nil
}

// Merge takes a config and combines it into a config file
func (k *KConf) Merge(config *clientcmdapi.Config, name string) {
	renamedClusters := make(map[string]string)
	renamedUsers := make(map[string]string)

	for clsName, cls := range config.Clusters {
		added := k.AddCluster(clsName, cls)
		if added != "" { // this cluster was newly added
			if added != clsName {
				Out.Log().Msgf("renamed cluster '%s' to '%s'", clsName, added)
				renamedClusters[clsName] = added
			}
			Out.Log().Msgf("added cluster '%s'", added)
		}
	}
	for uName, user := range config.AuthInfos {
		added := k.AddUser(uName, user)
		if added != "" { // this user was newly added
			if added != uName {
				Out.Log().Msgf("renamed user '%s' to '%s'", uName, added)
				renamedUsers[uName] = added
			}
			Out.Log().Msgf("added user '%s'", added)
		}
	}
	for ctxName, ctx := range config.Contexts {
		// if anything got renamed, we need to update this context
		if renamed, ok := renamedClusters[ctx.Cluster]; ok {
			ctx.Cluster = renamed
		}
		if renamed, ok := renamedUsers[ctx.AuthInfo]; ok {
			ctx.AuthInfo = renamed
		}
		// the context name is chosen based on the following priority:
		// 1. provided via --context-name
		// 2. specified in the provided config file
		// 3. generated from the cluster and user name (<cluster>-<user>)
		if name != "" {
			ctxName = name
		} else if ctxName == "" && name == "" {
			ctxName = fmt.Sprintf("%s-%s", ctx.Cluster, ctx.AuthInfo)
		}
		added := k.AddContext(ctxName, ctx)
		if added != "" { // this context was newly added
			if added != ctxName {
				Out.Log().Msgf("renamed context '%s' to '%s'", ctxName, added)
			}
			Out.Log().Msgf("added context '%s'", added)
		}
	}
}

// SetNamespace sets the namespace for the context `contextName`
func (k *KConf) SetNamespace(contextName, namespace string) error {
	if _, ok := k.Config.Contexts[contextName]; !ok {
		return fmt.Errorf("could not find context '%s'", contextName)
	}
	k.Config.Contexts[contextName].Namespace = namespace

	return nil
}

// SetCurrentContext sets the kubeconfig current context to `currentContext`
func (k *KConf) SetCurrentContext(contextName string) error {
	if _, ok := k.Config.Contexts[contextName]; !ok {
		return fmt.Errorf("could not find context '%s'", contextName)
	}
	k.Config.CurrentContext = contextName

	return nil
}
