package kubeconfig

import "fmt"

// Remove takes a context and its associated resources out of kubeconfig
func (k *KConf) Remove(name string) error {
	context, ok := k.Contexts[name]
	if !ok { // the context never existed
		return fmt.Errorf("could not find context '%s'", name)
	}

	users := k.contextsWithUser(context.AuthInfo)
	if len(users) == 1 { // it's just this context
		delete(k.AuthInfos, context.AuthInfo)
		Out.Log().Msgf("removed '%s' user", context.AuthInfo)
	}
	delete(k.Clusters, context.Cluster)
	Out.Log().Msgf("removed '%s' cluster", context.Cluster)
	delete(k.Contexts, name)
	Out.Log().Msgf("removed '%s' context", name)

	return k.Save()
}
