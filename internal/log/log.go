package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	prefixFormat       string      = "[%s: %s-%s] "
	logDirPathFormat   string      = "%s/log"
	logFileFormat      string      = "%s/%s.log"
	permissionsLogDir  os.FileMode = 0764
	permissionsLogFile os.FileMode = 0664
)

func InitLogConfiguration(
	projectDir string,
	env string,
	appName string,
	serviceName string,
	logOutputEnabled bool,
) error {
	log.SetFlags(log.LstdFlags)
	log.SetPrefix(fmt.Sprintf(prefixFormat, appName, serviceName, env))

	logDirPath := fmt.Sprintf(logDirPathFormat, projectDir)
	if err := os.MkdirAll(logDirPath, permissionsLogDir); err != nil {
		return err
	}

	logFile, err := os.OpenFile(
		fmt.Sprintf(logFileFormat, logDirPath, env),
		os.O_WRONLY|os.O_CREATE,
		permissionsLogFile,
	)
	if err != nil {
		return err
	}

	stdOutput := io.Discard
	if logOutputEnabled {
		stdOutput = os.Stdout
	}

	log.SetOutput(
		io.MultiWriter(logFile, stdOutput),
	)

	return nil
}
