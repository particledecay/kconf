package kubeconfig

import "fmt"

// Remove takes a context and its associated resources out of kubeconfig
func (k *KConf) Remove(name string) error {
	context, ok := k.Contexts[name]
	if !ok { // the context never existed
		return fmt.Errorf("could not find context '%s'", name)
	}

	Out.Log().Msgf("removing '%s' user", context.AuthInfo)
	delete(k.AuthInfos, context.AuthInfo)
	Out.Log().Msgf("removing '%s' cluster", context.Cluster)
	delete(k.Clusters, context.Cluster)
	Out.Log().Msgf("removing '%s' context", name)
	delete(k.Contexts, name)

	return k.Save()
}
