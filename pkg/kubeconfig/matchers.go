package kubeconfig

import (
	"fmt"

	"github.com/onsi/gomega/format"
	gomatchers "github.com/onsi/gomega/matchers"
	gotypes "github.com/onsi/gomega/types"
)

type ContainContextMatcher struct {
	Element interface{}
}

func ContainContext(element interface{}) gotypes.GomegaMatcher {
	return &ContainContextMatcher{
		Element: element,
	}
}

func (matcher *ContainContextMatcher) Match(actual interface{}) (bool, error) {
	config, ok := actual.(*KConf)
	if !ok {
		return false, fmt.Errorf("ContainContext matcher expects a KConf. Got:\n%s", format.Object(actual, 1))
	}

	elemMatcher := &gomatchers.EqualMatcher{Expected: matcher.Element}

	var lastError error
	contexts, _ := config.List()
	for _, ctx := range contexts {
		success, err := elemMatcher.Match(ctx)
		if err != nil {
			lastError = err
			continue
		}
		if success {
			return true, nil
		}
	}

	return false, lastError
}

func (matcher *ContainContextMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain context matching", matcher.Element)
}

func (matcher *ContainContextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain element matching", matcher.Element)
}
