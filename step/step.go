package step

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/cache"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/pathutil"
)

const (
	stepId = "save-gradle-cache"

	// Cache key template
	// OS + Arch: to guarantee that stack-specific content (absolute paths, binaries) are stored separately
	// checksum values:
	// - `**/*.gradle*`: Gradle build files in any submodule, including ones written in Kotlin (*.gradle.kts)
	// - `**/gradle-wrapper.properties`: contains exact Gradle version
	// - `gradle.properties`: contains Gradle config values
	// - `gradle/libs.versions.toml`: version catalog file, contains dependencies and their versions
	key = `{{ .OS }}-{{ .Arch }}-gradle-cache-{{ checksum "**/*.gradle*" "**/gradle-wrapper.properties" "gradle.properties" "gradle/libs.versions.toml" }}`
)

// Cached paths
var paths = []string{
	// Local Gradle cache folder (contains build cache too when `org.gradle.caching = true`)
	"~/.gradle/caches",

	// Cache of downloaded Gradle binary
	"~/.gradle/wrapper",

	// Configuration cache in project-level folder
	".gradle/configuration-cache",
}

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type SaveCacheStep struct {
	logger       log.Logger
	inputParser  stepconf.InputParser
	pathChecker  pathutil.PathChecker
	pathProvider pathutil.PathProvider
	pathModifier pathutil.PathModifier
	envRepo      env.Repository
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	pathChecker pathutil.PathChecker,
	pathProvider pathutil.PathProvider,
	pathModifier pathutil.PathModifier,
	envRepo env.Repository,
) SaveCacheStep {
	return SaveCacheStep{
		logger:       logger,
		inputParser:  inputParser,
		pathChecker:  pathChecker,
		pathProvider: pathProvider,
		pathModifier: pathModifier,
		envRepo:      envRepo,
	}
}

func (step SaveCacheStep) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf("failed to parse inputs: %w", err)
	}
	stepconf.Print(input)
	step.logger.Println()
	step.logger.Printf("Cache key: %s", key)
	step.logger.Printf("Cache paths:")
	step.logger.Printf(strings.Join(paths, "\n"))
	step.logger.Println()

	step.logger.EnableDebugLog(input.Verbose)

	saver := cache.NewSaver(step.envRepo, step.logger, step.pathProvider, step.pathModifier, step.pathChecker)
	return saver.Save(cache.SaveCacheInput{
		StepId:      stepId,
		Verbose:     input.Verbose,
		Key:         key,
		Paths:       paths,
		IsKeyUnique: true,
	})
}
