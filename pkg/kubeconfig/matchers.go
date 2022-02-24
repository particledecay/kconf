package kubeconfig

import (
	"fmt"

	"github.com/onsi/gomega/format"
	gomatchers "github.com/onsi/gomega/matchers"
	gotypes "github.com/onsi/gomega/types"
)

// ContainContextMatcher gets returned by the matcher function
type ContainContextMatcher struct {
	Element interface{}
}

// ContainContext is the Gomega matcher function
func ContainContext(element interface{}) gotypes.GomegaMatcher {
	return &ContainContextMatcher{
		Element: element,
	}
}

// Match looks in a kubeconfig for a matching Context
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

// FailureMessage displays when the context is not found but should be
func (matcher *ContainContextMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to contain context matching", matcher.Element)
}

// NegatedFailureMessage displays when the context is found but shouldn't be
func (matcher *ContainContextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to contain element matching", matcher.Element)
}
