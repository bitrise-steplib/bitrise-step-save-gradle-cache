package main

import (
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
	"github.com/bitrise-steplib/steps-save-gradle-cache/step"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)
	pathChecker := pathutil.NewPathChecker()
	pathProvider := pathutil.NewPathProvider()
	pathModifier := pathutil.NewPathModifier()
	cacheStep := step.New(logger, inputParser, pathChecker, pathProvider, pathModifier, envRepo)

	exitCode := 0

	if err := cacheStep.Run(); err != nil {
		logger.Errorf(err.Error())
		exitCode = 1
		return exitCode
	}

	return exitCode
}
