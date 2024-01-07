package repository_test

import (
	pconfig "github.com/raffaele-pilloni/axxon-test/config"
	clog "github.com/raffaele-pilloni/axxon-test/internal/log"
	"io"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRepository(t *testing.T) {
	config, err := pconfig.LoadConfig(true)
	if err != nil {
		log.Panicf("Load configuration failed. error: %v", err)
	}

	if err := clog.InitLogConfiguration(
		config.App.ProjectDir,
		config.App.Env,
		config.App.AppName,
		config.App.ServiceName,
		io.Discard,
	); err != nil {
		log.Panicf("Init log configuration failed. error: %v", err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
