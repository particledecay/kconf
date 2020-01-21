package kubeconfig

import (
	"fmt"
)

// Remove takes a context and its associated resources out of kubeconfig
func (k *KConf) Remove(name string) error {
	context, ok := k.Contexts[name]
	if !ok { // the context never existed
		return fmt.Errorf("Could not find context '%s'", name)
	}

	fmt.Printf("Removing %s user\n", context.AuthInfo)
	delete(k.AuthInfos, context.AuthInfo)
	fmt.Printf("Removing %s cluster\n", context.Cluster)
	delete(k.Clusters, context.Cluster)
	fmt.Printf("Removing %s context\n", name)
	delete(k.Contexts, name)

	return k.Save()
}
