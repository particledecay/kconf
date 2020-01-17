package kubeconfig

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"k8s.io/client-go/tools/clientcmd"
)

// Write takes a config and saves it to a file
func Write(config *KConf) error {
	err := clientcmd.WriteToFile(config.Config, MainConfigPath)
	if err != nil {
		log.Error().Msgf("Could not write to file at %s", MainConfigPath)
	}
	return nil
}

// Merge takes a config and combines it into a config file
func Merge(config *KConf) error {
	mainConfig, err := Read(MainConfigPath)
	if err != nil {
		log.Fatal().Msgf("Could not read main config file: %v", err)
	}

	log.Printf("%v", mainConfig)

	// merge users
	for user := range config.AuthInfos {
		if _, ok := mainConfig.AuthInfos[user]; !ok {
			mainConfig.AuthInfos[user] = config.AuthInfos[user]
			fmt.Printf("Adding new user %s", user)
		}
	}

	// merge clusters
	for cluster := range config.Clusters {
		if _, ok := mainConfig.Clusters[cluster]; !ok {
			mainConfig.Clusters[cluster] = config.Clusters[cluster]
			fmt.Printf("Adding new cluster %s", cluster)
		}
	}

	// merge contexts
	for ctxName, ctx := range config.Contexts {
		if added := mainConfig.AddContext(ctx); added != "" {
			fmt.Printf("Adding new context '%s'", added)
		}
		fmt.Printf("Adding new context '%s'", ctxName)
		mainConfig.Contexts[ctxName] = ctx
	}

	err = Write(mainConfig)
	if err != nil {
		log.Fatal().Msgf("Error while writing merged config: %v", err)
	}

	return nil
}
