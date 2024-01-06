package entity_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIncomingEvent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Entity Suite")
}
